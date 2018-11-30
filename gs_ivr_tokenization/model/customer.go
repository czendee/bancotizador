package model
import (
    "database/sql"
    "errors"
	"fmt"
	"log"
	"strings"
)
type Customer struct {
    ID    string    `sql:"type:bigserial"`
    Reference   string    `sql:"type:varchar(100)`
    Createat   string    `sql:"type:timestamp`
    Lastupdateat   string    `sql:"type:timestamp`
}

func (u *Customer) GetCustomer(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT id_customer, reference FROM banwirecustomer WHERE id_customer=%d", u.ID)
    return db.QueryRow(statement).Scan(&u.ID, &u.Reference)
}

func (u *Customer) GetCustomerByReference01(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT id_customer FROM banwirecustomer WHERE reference='%s'", u.Reference)
    errorselect :=db.QueryRow(statement).Scan(&u.ID)
    log.Print("ejecuto el select")
    if errorselect != nil {
        log.Print("ejecuto el select: error :"+errorselect.Error())
        var textoerr =errorselect.Error()
        
        if strings.Contains(textoerr,"no rows in result set"){
        	 log.Print("no existe, entocnes inserta customer :")
	//	    statement := fmt.Sprintf("INSERT INTO app(id_customer, reference, created_at, last_modified_at) VALUES('7777','%s',current_timestamp,current_timestamp)", u.ID,  u.Reference)
		    statement := fmt.Sprintf("INSERT INTO banwirecustomer( reference, created_at, last_update_at) VALUES('%s',current_timestamp,current_timestamp)",   u.Reference)
		    _, errExec := db.Exec(statement)
		    if errExec != nil {
        	    log.Print("inserta customer y hay algo")
                var textoerrexec =errExec.Error()
		    	if strings.Contains(textoerrexec,"no rows in result set"){
				    statement := fmt.Sprintf("SELECT id_customer FROM banwirecustomer WHERE reference='%s'", u.Reference)
				    errorselect02 :=db.QueryRow(statement).Scan(&u.ID)
				    log.Print("ejecuto el select 2")
				    if errorselect02 != nil {
				    	//no existia  el customer,  insert e intento dar select al registor nuevo, pero ocurrio error
				    	return errorselect02
	                }else{
				    	//no existia  el customer,  insert y le dio select al registor nuevo, Exito, regresa el ID_customer	                	
				    	return nil
	                }
		    		
                }else{
		            log.Print("ejecuto el select: error Insert :"+errExec.Error())
		            errorselect=errExec                	
		            
                }  

		    }else{
				   //no existia  el customer, intento insert y ok
        	    log.Print("inserta customer y ahora consulta el id_customer")

				    statement := fmt.Sprintf("SELECT id_customer FROM banwirecustomer WHERE reference='%s'", u.Reference)
				    errorselect02 :=db.QueryRow(statement).Scan(&u.ID)
				    log.Print("ejecuto el select 2")
				    if errorselect02 != nil {
				    	//no existia  el customer,  insert e intento dar select al registor nuevo, pero ocurrio error
				    	return errorselect02
	                }else{
				    	//no existia  el customer,  insert y le dio select al registor nuevo, Exito, regresa el ID_customer	                	
				    	return nil
	                }


		    }
        	
        }

    }else{
     //errorSelect is nil
	//ya exoiste el customer y regresa el id
    }
   return errorselect
//    return db.QueryRow(statement).Scan(&u.ID, &u.Reference)
}

func (u *Customer) updateCustomer(db *sql.DB) error {
    statement := fmt.Sprintf("UPDATE banwirecustomer SET reference='%s', last_update_at= current_timestamp WHERE id_customer=%s", u.Reference,  u.ID)
    _, err := db.Exec(statement)
    return err
}
func (u *Customer) deleteCustomer(db *sql.DB) error {
    return errors.New("Not implemented")
}
func (u *Customer) hazCreateCustomer(db *sql.DB) error {

    statement := fmt.Sprintf("INSERT INTO banwirecustomer(id_customer, reference, created_at, last_modified_at) VALUES(%s,'%s',current_timestamp,current_timestamp)", u.ID,  u.Reference)
    _, err := db.Exec(statement)
    if err != nil {
        return err
    }
    err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)
    if err != nil {
        return err
    }
    return nil
}
func GetCustomerByReference(db *sql.DB, reference string) (Customer, error) {
	  log.Print("GetCustomerByReference 01!\n")
	  var resultado Customer
    statement := fmt.Sprintf("SELECT id_customer,  reference FROM banwirecustomer WHERE reference=%d", reference)
    //return
    db.QueryRow(statement).Scan(&resultado.ID, &resultado.Reference)
    return resultado,nil
}

func getCustomers(db *sql.DB, start, count int) ([]Customer, error) {
 statement := fmt.Sprintf("SELECT id_customer,  reference FROM banwirecustomer LIMIT %d OFFSET %d", count, start)
    rows, err := db.Query(statement)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    cards := []Customer{}
    for rows.Next() {
        var u Customer
        if err := rows.Scan(&u.ID, &u.Reference); err != nil {
            return nil, err
        }
        cards = append(cards, u)
    }
    return cards, nil
}