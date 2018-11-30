package db

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	pg "github.com/jackc/pgx"
)

// PgDbName is name of postgresql database
const PgDbName DbName = "PgDb"

// PgDb represents the postgresql database
type PgDb struct {
	loaded bool
	pool   *pg.ConnPool

	mutex *sync.Mutex
}

// PgDbConn handles and stores the database connection
type PgDbConn struct {
	*pg.Conn
	pool *pg.ConnPool
}

// Get returns the database connection stored
func (d *PgDbConn) Get() interface{} {
	return d.Conn
}

// Close closes the database connection stored
func (d *PgDbConn) Close() error {
	d.pool.Release(d.Conn)
	return nil
}

// GetName returns the database name
func (d *PgDb) GetName() DbName {
	return PgDbName
}

// NewPgDb returns a new postgresql database object (*PgDb)
func NewPgDb(host string, port int, db, user, password string) *PgDb {
	config := pg.ConnPoolConfig{ConnConfig: pg.ConnConfig{
		Host:     host,
		Port:     uint16(port),
		User:     user,
		Password: password,
		Database: db,
	}, MaxConnections: 20}

	var _db = &PgDb{}

	var sqldb, err = pg.NewConnPool(config)

	if err == nil {
		_db.pool = sqldb
		_db.mutex = new(sync.Mutex)
		_db.loaded = true
	} else {
		log.Println(err)
	}

	return _db

}

// Loaded returns if the database was loaded
func (d *PgDb) Loaded() bool {
	return d.loaded
}

// IsMaster returns if the connection is master
func (d *PgDb) IsMaster() bool {
	var isMaster bool

	var result string
	if err := d.pool.QueryRow("SELECT pg_is_in_recovery()").Scan(&result); err == nil {
		isMaster = (result == "f" || result == "false")
	}

	return isMaster
}

// Status returns if the connection is available
func (d *PgDb) Status() bool {
	var ok bool

	var result string
	if err := d.pool.QueryRow("SELECT pg_is_in_recovery()").Scan(&result); err == nil {
		ok = true
	}

	return ok
}

// Get returns a new connection of the database
func (d *PgDb) Get() dbConn {

	if _db, err := d.pool.Acquire(); err == nil {
		return &PgDbConn{_db, d.pool}
	}
	return nil
}

// GetPgDb returns the PgDb connection associated to the reference
func (d *dbs) GetPgDb(reference interface{}, master bool) *pg.Conn {
	if _d, ok := d.Get(reference, PgDbName, master).(*pg.Conn); ok {
		return _d
	}

	return nil
}

//PgLAO handels serialization of SQL params object in string
func PgLAO(args ...interface{}) string {
	var strs = []string{}
	for _, a := range args {
		switch v := a.(type) {
		case *string, *int8, *int16, *int32, *int64, *float32, *float64, *bool:

			strs = append(strs, fmt.Sprintf("%#v", reflect.ValueOf(v).Elem()))
			break
		case *sql.NullString, *sql.NullInt64, *sql.NullFloat64, *sql.NullBool:
			if val, _ := v.(sqlValue).Value(); val != nil {
				strs = append(strs, fmt.Sprintf("%#v", val))
			} else {
				strs = append(strs, "null")
			}

			break
		case *DefaultString, *DefaultInt64, *DefaultFloat64, *DefaultBool:
			var df = v.(sqlValueDefault)
			if val, _ := df.Value(); val != nil {
				strs = append(strs, fmt.Sprintf("%#v", val))
			} else {
				strs = append(strs, fmt.Sprintf("%v", df.Default()))
			}
			break
		case *Timestamp:
			if v.IsZero() {
				strs = append(strs, "")
			} else if s := v.Format("2006-01-02T15:04:05.999999"); s != "" {
				strs = append(strs, fmt.Sprintf("%#v", s))
			}
		case *Date:
			if v.IsZero() {
				strs = append(strs, "")
			} else if s := v.Format("2006-01-02"); s != "" {
				strs = append(strs, fmt.Sprintf("%#v", s))
			}
			break
		case *Time:
			if v.IsZero() {
				strs = append(strs, "")
			} else if s := v.Format("15:04:05.999999-07"); s != "" {
				strs = append(strs, fmt.Sprintf("%#v", s))
			}
			break
		case *PgJson:
			if val, _ := v.Value(); val != nil {
				strs = append(strs, fmt.Sprintf("%#v", val))
			} else {
				strs = append(strs, "null")
			}
			break
		case *ArrayString, *ArrayInt64:
			var df = v.(sqlValue)
			if val, _ := df.Value(); val != nil {
				strs = append(strs, fmt.Sprintf("%#v", ToPgArray(val)))
			} else {
				strs = append(strs, "{}")
			}
			break
		default:
			var _value = reflect.ValueOf(v)
			var _type = _value.Type()

			if _value.Kind() == reflect.Ptr {
				_value = _value.Elem()
				_type = _type.Elem()
				v = _value.Interface()
			}

			switch _type.Kind() {
			case reflect.Slice:
				strs = append(strs, fmt.Sprintf("%#v", ToPgArray(v)))
				break
			default:
				log.Println(v, reflect.TypeOf(v))
				break
			}
			break
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(strs, ","))
}

// ToPgArray converts a array or slice object to string for SQL postgresql
func ToPgArray(t interface{}) string {
	var slice = []string{}
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			switch v := s.Index(i).Interface().(type) {
			case string:
				var escaper = strings.NewReplacer(
					`\`, `\\`,
					`'`, `\'`,
				)
				slice = append(slice, fmt.Sprintf("%v", escaper.Replace(v)))
				break
			default:
				switch reflect.TypeOf(v).Kind() {
				case reflect.Slice:
					slice = append(slice, ToPgArray(v))
					break
				default:
					slice = append(slice, fmt.Sprintf("%v", v))
					break
				}
				break
			}

		}
	}

	return fmt.Sprintf("{%s}", strings.Join(slice, ","))
}

// ToPgJson converts any object to string JSON for SQL postgresql
func ToPgJson(t interface{}) interface{} {
	if b, e := json.Marshal(t); e == nil {
		if s := string(b); s == "null" {
			return nil
		} else {
			return s
		}
	}

	return nil
}

// PgArrayToArray converts postgresql array string to a slice string
func PgArrayToArray(array string) ([]string, error) {
	var out []string
	var arrayOpened, quoteOpened, escapeOpened bool
	item := &bytes.Buffer{}
	for _, r := range array {
		switch {
		case !arrayOpened:
			if r != '{' {
				return nil, errors.New("Doesn't appear to be a postgres array.  Doesn't start with an opening curly brace.")
			}
			arrayOpened = true
		case escapeOpened:
			item.WriteRune(r)
			escapeOpened = false
		case quoteOpened:
			switch r {
			case '\\':
				escapeOpened = true
			case '"':
				quoteOpened = false
				if item.String() == "NULL" {
					item.Reset()
				}
			default:
				item.WriteRune(r)
			}
		case r == '}':
			// done
			out = append(out, item.String())
			return out, nil
		case r == '"':
			quoteOpened = true
		case r == ',':
			// end of item
			out = append(out, item.String())
			item.Reset()
		default:
			item.WriteRune(r)
		}
	}
	return nil, errors.New("Doesn't appear to be a postgres array.  Premature end of string.")
}

// ArrayString is a object that support scan and value functions for array object postgresql
type ArrayString struct {
	Elms  []string
	Valid bool
}

// Scan implements the Scanner interface
func (a *ArrayString) Scan(value interface{}) error {
	a.Valid = false
	switch v := value.(type) {
	case *[]string:
		a.Elms = *v
		a.Valid = true
		break
	case []string:
		a.Elms = v
		a.Valid = true
		break
	case string:
		a.Elms, _ = PgArrayToArray(v)
		a.Valid = true
		break
	}
	return nil
}

// Value implements the driver Valuer interface
func (a *ArrayString) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Elms, nil
}

// ArrayInt64 is a object that support scan and value functions for array object postgresql
type ArrayInt64 struct {
	Elms  []int64
	Valid bool
}

// Scan implements the Scanner interface
func (a *ArrayInt64) Scan(value interface{}) error {
	a.Valid = false
	switch v := value.(type) {
	case *[]int64:
		a.Elms = *v
		a.Valid = true
		break
	case []int64:
		a.Elms = v
		a.Valid = true
		break
	case string:
		a.Elms = []int64{}
		_a, _ := PgArrayToArray(v)
		for _, s := range _a {
			if i, e := strconv.ParseInt(s, 10, 64); e == nil {
				a.Elms = append(a.Elms, i)
			}
		}
		a.Valid = true
		break
	}
	return nil
}

// Value implements the driver Valuer interface
func (a *ArrayInt64) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Elms, nil
}

// PgTable returns a new PgTableStruct object
func PgTable(tbl interface{}) *PgTableStruct {
	return &PgTableStruct{
		tbl: tbl,
	}
}

// PgTableStruct is a object that supports the table postgresql object methods
type PgTableStruct struct {
	tbl interface{}
	mdl interface{}
}

// Close copys the table values to model
func (tblStruct *PgTableStruct) Close() {
	copyTo(tblStruct.mdl, tblStruct.tbl, "namefield")
}

// Model sets the model to the PgTableStruct object
func (tblStruct *PgTableStruct) Model(mdl interface{}) *PgTableStruct {
	tblStruct.mdl = mdl
	copyTo(tblStruct.tbl, tblStruct.mdl, "namefield")
	return tblStruct
}

// Value implements the driver Valuer interface
func (tblStruct *PgTableStruct) Value() (driver.Value, error) {
	return PgLAO(tblStruct.ToScan()...), nil
}

// ToScan returns and prepairs the paramters for be scanned
func (tblStruct *PgTableStruct) ToScan() []interface{} {
	var fields = []interface{}{}
	var nums = []int{}

	v := reflect.ValueOf(tblStruct.tbl).Elem()

	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		t := v.Type().Field(j)

		tag := t.Tag.Get("dbfield")
		stag := strings.Split(tag, ",")
		if len(stag) > 0 {
			if num, e := strconv.Atoi(stag[0]); e == nil {
				fields = append(fields, f.Addr().Interface())
				nums = append(nums, num)
			}
		}
	}

	sort.Sort(pgTableSort{fields, nums})

	return fields
}

// pgTableSort handles the sorting of an object
type pgTableSort struct {
	obj  []interface{}
	nums []int
}

// Len is the number of elements in the collection
func (a pgTableSort) Len() int { return len(a.obj) }

// Swap swaps the elements with indexes i and j
func (a pgTableSort) Swap(i, j int) {
	a.obj[i], a.obj[j] = a.obj[j], a.obj[i]
	a.nums[i], a.nums[j] = a.nums[j], a.nums[i]
}

// Less reports whether the element with
// index i should sort before the element with index j
func (a pgTableSort) Less(i, j int) bool { return a.nums[i] < a.nums[j] }

//PgParams returns the params in a slice
func PgParams(params ...interface{}) []interface{} {

	//fix for pgx driver
	for i, p := range params {
		if v, ok := p.(driver.Valuer); ok {
			params[i], _ = v.Value()
		}
	}

	return params
}

// PgJson is an object that supports json value for postgresql
type PgJson struct {
	Data  interface{}
	Valid bool
}

// Scan implements the Scanner interface
func (j *PgJson) Scan(value interface{}) error {
	if value == nil {
		j.Data, j.Valid = "", false
		return nil
	}
	j.Valid = true
	j.Data = value
	return nil
}

// Value implements the driver Valuer interface
func (j *PgJson) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return ToPgJson(j.Data), nil
}

// To handles the serialization and sets the value into the object
func (j *PgJson) To(i interface{}) error {
	var b []byte

	switch d := j.Data.(type) {
	case string:
		b = []byte(d)
		break
	case map[string]interface{}:
		b, _ = json.Marshal(d)
		break
	}

	return json.Unmarshal(b, i)
}

// MarshalJSON returns the string json of data object
func (j PgJson) MarshalJSON() ([]byte, error) {
	fmt.Println(j.Data)

	//fix for pgx driver
	switch v := j.Data.(type) {
	case string:
		return []byte(v), nil
		break
	}

	return json.Marshal(j.Data)
}

func PgsqlFetchFirst(c *pg.Tx, name string, args ...interface{}) bool {
	if err := c.QueryRow(fmt.Sprintf("FETCH FIRST FROM %s", name)).Scan(args...); err != nil {
		if err == pg.ErrNoRows {
			return false
		}

		panic(err)
	}

	return true
}

func PgSqlFetchNext(c *pg.Tx, name string, args ...interface{}) bool {
	if err := c.QueryRow(fmt.Sprintf("FETCH NEXT FROM %s", name)).Scan(args...); err != nil {
		if err == pg.ErrNoRows {
			return false
		}

		panic(err)
	}

	return true
}

// PgSqlFetch handelds the fetch process for postgresql
func PgSqlFetch(c *pg.Tx, ch chan bool, name string, args ...interface{}) {
	if rws, err := c.Query(fmt.Sprintf("FETCH ALL FROM %s", name)); err == nil {
		defer func() {
			rws.Close()
		}()

		if err = rws.Err(); err != nil {
			panic(err)
		}

		for rws.Next() {
			if err = rws.Scan(args...); err != nil {
				panic(err)
			}

			if ch != nil {
				ch <- true
				_ = <-ch
			}
		}

		if ch != nil {
			ch <- false
		}

		if err = rws.Err(); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}

// PgSqlFetchRow handelds the fetch process for postgresql (fetch one)
func PgSqlFetchRow(c *pg.Tx, dir string, name string, args ...interface{}) {
	if dir == "" {
		dir = "NEXT"
	}
	rw := c.QueryRow(fmt.Sprintf("FETCH %s FROM %s", dir, name))
	if rw != nil {
		if err := rw.Scan(args...); err != nil {
			panic(err)
		}
	}
}
