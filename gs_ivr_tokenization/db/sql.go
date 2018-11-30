package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// defaultType stores the default value
type defaultType struct {
	val interface{}
}

// Default returns the default value stored
func (t *defaultType) Default() interface{} {
	if t.val == nil {
		return ""
	}
	return t.val
}

// DefaultString stores the default value of string type
type DefaultString struct {
	sql.NullString
	defaultType
}

// NewDefaultString returns a new DefaultString object
func NewDefaultString(v string, b bool, d interface{}) DefaultString {
	return DefaultString{
		sql.NullString{v, b},
		defaultType{d},
	}
}

// DefaultInt64 stores the default value of int64 type
type DefaultInt64 struct {
	sql.NullInt64
	defaultType
}

// NewDefaultInt64 returns a new DefaultInt64 object
func NewDefaultInt64(v int64, b bool, d interface{}) DefaultInt64 {
	return DefaultInt64{
		sql.NullInt64{v, b},
		defaultType{d},
	}
}

// DefaultFloat64 stores the default value of float64 type
type DefaultFloat64 struct {
	sql.NullFloat64
	defaultType
}

// NewDefaultFloat64 returns a new DefaultFloat64 object
func NewDefaultFloat64(v float64, b bool, d interface{}) DefaultFloat64 {
	return DefaultFloat64{
		sql.NullFloat64{v, b},
		defaultType{d},
	}
}

// DefaultBool stores the default value of boolean type
type DefaultBool struct {
	sql.NullBool
	defaultType
}

// NewDefaultBool returns a new DefaultBool object
func NewDefaultBool(v bool, b bool, d interface{}) DefaultBool {
	return DefaultBool{
		sql.NullBool{v, b},
		defaultType{d},
	}
}

// Time is an object that supports time value for sql
type Time struct {
	time.Time
	Valid bool
}

// Scan implements the Scanner interface
func (nt *Time) Scan(value interface{}) error {
	nt.Valid = false

	switch t := value.(type) {
	case nil:
		nt.Time = time.Time{}
		nt.Valid = true
		break
	case time.Time:
		nt.Time = t
		nt.Valid = true
		break
	case string:
		var e error
		//"15:04:05.999999-07"
		if nt.Time, e = time.Parse("2006-01-02 15:04:05.999999", t); e == nil {
			nt.Valid = true
		} else if nt.Time, e = time.Parse("15:04:05.999999", t); e == nil {
			nt.Valid = true
		} else if nt.Time, e = time.Parse("15:04:05.999999-07", t); e == nil {
			nt.Valid = true
		} else if nt.Time, e = time.Parse("15:04:05.999999-07:00", t); e == nil {
			nt.Valid = true
		} else if nt.Time, e = time.Parse("15:04:05.999999-07:00:00", t); e == nil {
			nt.Valid = true
		}
		break
	default:
		if v, ok := t.(driver.Valuer); ok {
			if _v, e := v.Value(); e == nil {
				nt.Time = _v.(time.Time)
				nt.Valid = true
			}
		}

		break
	}
	return nil
}

// Value implements the driver Valuer interface
func (nt *Time) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// SetParse handles the parsing of a string value to time object
func (nt *Time) SetParse(s string, f string) {
	var e error
	nt.Valid = false
	if nt.Time, e = time.Parse(f, s); e == nil {
		nt.Valid = true
	}
}

// Date is an object that supports date value for sql
type Date struct {
	Time
}

// NewDateByString converts string to date object and returns
func NewDateByString(s string, f string) (d Date) {
	d = Date{}
	d.SetParse(s, f)
	return
}

// Timestamp is an object that supports timestamp value for sql
type Timestamp struct {
	Time
}

// sqlValue is an interface that suports the sql value object methods
type sqlValue interface {
	Value() (driver.Value, error)
}

// sqlValueDefault is an interface that supports the sql value default object methods
type sqlValueDefault interface {
	Value() (driver.Value, error)
	Default() interface{}
}

// SqlFetch handelds the fetch process for sql
func SqlFetch(c *sql.Tx, ch chan bool, name string, args ...interface{}) {
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
