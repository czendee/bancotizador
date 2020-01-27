package model
import (
    "database/sql"
//    "errors"
	"fmt"
	"log"
//	miu "banwire/services/gs_ivr_tokenization/util"
	
)
type Payment struct {
    ID    string    `sql:"type:bigserial"`
    Token  string `sql:"type:varchar(30)`
    Created_at   string    `sql:"type:timestamp`
    Amount  string    `sql:"type:bigint`
    Msi  string    `sql:"type:integer`
}
func (u *Payment) getPayment(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT token, created_at FROM banwirepayment WHERE id=%d", u.ID)
    return db.QueryRow(statement).Scan(&u.Token, &u.Created_at)
}
func GetTodayPaymentsByTokenCard(db *sql.DB, eltoken string ) ([]Payment, error) {
//     dateDDMMYY, errorGeneral := miu.ConvertMMYYintoDDMMYY(u.CreatedValid)  //utils.go
//     if(errorGeneral==""){
//     	log.Print("fecha valid"+dateDDMMYY)
     	log.Print("procesando GetTodayPaymentsByTokenCard")
//        statement := fmt.Sprintf("SELECT token,created_at,amount FROM banwirepayment WHERE token='%s' and created_at  AT TIME ZONE 'cst' > date_trunc('day', now()) AT TIME ZONE 'cst' ",eltoken)
//de aqui se saco como checar las payments del dia
//se usa el inteval -6h para mexico
//SELECT token,created_at,created_at- interval '6h' as mera, now()::date + interval '1minutes'  as pasadamedia, amount ,created_at  AT TIME ZONE 'cst', now() as ahora,date_trunc('day', now()) AT TIME ZONE 'cst' as ui, now()::date as unn, now()::date + interval '1h' as hha
// FROM banwirepayment WHERE token='crd.1jJckdSrFY3uAVZecM1voax1gLhv' 
//--and created_at < now()::date + interval '1h'
//--and created_at- interval '6h' >=  now()::date + interval '2minutes' 

        statement := fmt.Sprintf("SELECT token,created_at,amount FROM banwirepayment WHERE token='%s' and created_at- interval '6h' >= (now()- interval '6h')::date + interval '1minutes' ",eltoken)
        
//        return db.QueryRow(statement).Scan(&u.Token, &u.Created_at,&u.Amount),nil

log.Print("procesando GetTodayPaymentsByTokenCard"+"SELECT token,created_at,amount FROM banwirepayment WHERE token='%s' and created_at- interval '6h' >= (now()- interval '6h')::date + interval '1minutes' ",eltoken)

    rows, err := db.Query(statement)
    log.Print("GetTodayPaymentsByTokenCard 02.1!\n")
    if err != nil {
        return nil, err
    }
    log.Print("GetTodayPaymentsByTokenCard 02.5!\n")
    defer rows.Close()
    dailypayments := []Payment{}
    for rows.Next() {
    	 log.Print("GetTodayPaymentsByTokenCard 03!\n")
        var u Payment
        if err := rows.Scan( &u.Token, &u.Created_at,&u.Amount); err != nil {

        	log.Print("GetTodayPaymentsByTokenCard err!\n"+err.Error())
            return nil, err
        }
    	 log.Print("GetTodayPaymentsByTokenCard 03.5!\n")
        dailypayments = append(dailypayments, u)
    }
    log.Print("GetTodayPaymentsByTokenCard 04!\n")
    return dailypayments, nil

  //   }else{
  //       return errorGeneral
  //   }
}

func GetAllPaymentsByTokenCard(db *sql.DB, eltoken string ) (string, error) {
     	log.Print("procesando GetTodayPaymentsByTokenCard")
        statement := fmt.Sprintf("SELECT token,created_at,amount FROM banwirepayment WHERE token='%s'  ",eltoken)
        
//        return db.QueryRow(statement).Scan(&u.Token, &u.Created_at,&u.Amount),nil

log.Print("procesando GetTodayPaymentsByTokenCard"+"SELECT token,created_at,amount FROM banwirepayment WHERE token='%s'  ",eltoken)

    rows, err := db.Query(statement)
    log.Print("GetAllPaymentsByTokenCard 02.1!\n")
    if err != nil {
        return "", err
    }
    
    log.Print("GetAllPaymentsByTokenCard 02.5!\n")
    defer rows.Close()
    
    var cuantos int
    cuantos=0
    for rows.Next() {
    	 log.Print("GetAllPaymentsByTokenCard 03!\n")
         cuantos=cuantos+1
    }
    log.Print("GetAllPaymentsByTokenCard 04!\n")
    if(cuantos>0){
        log.Print("GetAllPaymentsByTokenCard 05! existen previos\n")
         return "EXISTEN PREVIOS", nil;    
    }else{
        log.Print("GetAllPaymentsByTokenCard 05! no existen previos\n")
         return "No EXISTEN PREVIOS", nil;      
    }
  
}

func (u *Payment) CreatePayment(db *sql.DB) error {

	    statement := fmt.Sprintf("INSERT INTO banwirepayment( token, created_at, amount, msi) VALUES('%s',current_timestamp,%s,%s)",  u.Token, u.Amount, u.Msi)
	    _, err := db.Exec(statement)
     	log.Print("exec ejecutado")
	    if err != nil {
     	   log.Print("exec ejecutado:error "+err.Error())
	        return err
	    }
     	   log.Print("exec ejecutado:todo ok ")
    return nil
}
func GetPaymentsCardsByCustomer(db *sql.DB, customer_reference string) ([]Card, error) {
	log.Print("GetCardsByCustomer 01!\n")
//	statement := fmt.Sprintf("SELECT id_card, token, bin,last_digits,valid_thru,brand,type_card, score FROM card WHERE id_customer= %s  order by score DESC", id_customer)
	statement := fmt.Sprintf("SELECT id_card, token, bin,last_digits,valid_thru,brand,type_card, score  FROM banwirecard  c, banwirecustomer  a WHERE c.id_customer= a.id_customer and a.reference= '%s'  order by score DESC;", customer_reference)
 	log.Print("GetCardsByCustomerb 02!\n")
    rows, err := db.Query(statement)
    log.Print("GetCardsByCustomer 02.1!\n")
    if err != nil {
        return nil, err
    }
    log.Print("GetCardsByCustomer 02.5!\n")
    defer rows.Close()
    cards := []Card{}
    for rows.Next() {
    	 log.Print("GetCardsByCustomer 03!\n")
        var u Card
        if err := rows.Scan(&u.ID, &u.Token, &u.Bin,&u.Last,&u.Valid, &u.Brand,&u.Type, &u.Score); err != nil {
//        if err := rows.Scan(&u.Token, &u.Bin); err != nil {
        	log.Print("GetCardsByCustomer err!\n"+err.Error())
            return nil, err
        }
    	 log.Print("GetCardsByCustomer 03.5!\n")
        cards = append(cards, u)
    }
    log.Print("GetCardsByCustomer 04!\n")
    return cards, nil
}
