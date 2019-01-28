package model
import (
    "database/sql"
    "errors"
	"fmt"
	"log"
	miu "banwire/services/gs_ivr_tokenization/util"
	
)
type Card struct {
    ID    string    `sql:"type:bigserial"`
    Token  string `sql:"type:varchar(30)`
    Last   string    `sql:"type:varchar(4)`
    Bin   string    `sql:"type:varchar(6)`
    Valid   string    `sql:"type:timestamp`
    Score   string    `sql:"type:integer`
    Customer   string    `sql:"type:bigint`
    Brand   string    `sql:"type:varchar(30)`
    Type   string    `sql:"type:varchar(30)`
}
/////////////////////////////version 3.2 :)
func (u *Card) getCard(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT token, last FROM card WHERE id=%d", u.ID)
    return db.QueryRow(statement).Scan(&u.Token, &u.Last)
}
func (u *Card) GetCardByToken(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT token,last_digits,bin, brand, type_card  reference FROM banwirecard WHERE token='%s'", u.Token)
    return db.QueryRow(statement).Scan(&u.Token, &u.Last,&u.Bin,&u.Brand,&u.Type)
}
func (u *Card) GetCardByTokenAndCust(db *sql.DB,customer_id string) error {
    statement := fmt.Sprintf("SELECT token,last_digits,bin, brand, type_card  reference FROM banwirecard WHERE token='%s' and id_customer='%s' ", u.Token,customer_id)
    return db.QueryRow(statement).Scan(&u.Token, &u.Last,&u.Bin,&u.Brand,&u.Type)
}
/*
func GetCardByToken(db *sql.DB) (Customer, error) {
	  log.Print("GetCardByToken 01!\n")
	  var resultado Customer
    statement := fmt.Sprintf("SELECT token,last,  reference FROM card WHERE token=%d", u.Token)
    //return
    db.QueryRow(statement).Scan(&resultado.ID, &resultado.Reference)
    return resultado,nil
}
*/

func (u *Card) UpdateCard(db *sql.DB) error {
//    statement := fmt.Sprintf("UPDATE card SET score='%s', last_update_at= current_timestamp WHERE token='%s'", u.Score,  u.Token)
    statement := fmt.Sprintf("UPDATE banwirecard SET score=score +1, last_update_at= current_timestamp WHERE token='%s'",  u.Token)
    _, err := db.Exec(statement)
    return err
}

func (u *Card) IncreaseScoreCardAndCust(db *sql.DB,customer_id string) error {
//    statement := fmt.Sprintf("UPDATE card SET score='%s', last_update_at= current_timestamp WHERE token='%s'", u.Score,  u.Token)
    statement := fmt.Sprintf("UPDATE banwirecard SET score=score +1, last_update_at= current_timestamp WHERE token='%s' and id_customer='%s' ", u.Token,customer_id)
    _, err := db.Exec(statement)
    return err
}

func (u *Card) IncreaseScoreCard(db *sql.DB) error {
//    statement := fmt.Sprintf("UPDATE card SET score='%s', last_update_at= current_timestamp WHERE token='%s'", u.Score,  u.Token)
    statement := fmt.Sprintf("UPDATE banwirecard SET score=score +1, last_update_at= current_timestamp WHERE token='%s'",  u.Token)
    _, err := db.Exec(statement)
    return err
}

//Jan 27,2018: used in new rule to Remove card if 1st payment fails
func (u *Card) DeleteCard(db *sql.DB) error {
    statement := fmt.Sprintf("DELETE FROM banwirecard WHERE token='%s'",  u.Token)
    _, err := db.Exec(statement)
    return err
//    return errors.New("Not implemented")        
}

func (u *Card) deleteCard(db *sql.DB) error {
    return errors.New("Not implemented")
}
func (u *Card) CreateCard(db *sql.DB) error {
//    statement := fmt.Sprintf("INSERT INTO card(id_card, token, last_digits, bin, valid_thru,score, is_banned, banned_at,created_at, last_update_at, id_customer) VALUES(%d,'%s','%s','%s', %d, %d)", u.ID, u.Token, u.Last, u.Bin, u.score,u.Customer)

     dateDDMMYY, errorGeneral := miu.ConvertMMYYintoDDMMYY(u.Valid)  //utils.go
     if(errorGeneral==""){
     	log.Print("fecha valid"+dateDDMMYY)

	    statement := fmt.Sprintf("INSERT INTO banwirecard( token, last_digits, bin, valid_thru, score, is_banned, id_customer,brand, created_at, last_update_at, type_card) VALUES('%s','%s','%s',to_timestamp('%s 00:00:01', 'DD/MM/YYYY hh24:mi:ss'), %s,false, %s, '%s',current_timestamp,current_timestamp,'%s')",  u.Token, u.Last, u.Bin,dateDDMMYY, u.Score,u.Customer, u.Brand,u.Type)
	    _, err := db.Exec(statement)
     	log.Print("exec ejecutado")
	    if err != nil {
     	   log.Print("exec ejecutado:error "+err.Error())
	        return err
	    }
     	   log.Print("exec ejecutado:todo ok ")
/*	    
	    err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)
	    if err != nil {
	        return err
	    }
 */    	
     }
    return nil
}
func GetCardsByCustomer(db *sql.DB, customer_reference string) ([]Card, error) {
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

func getCards(db *sql.DB, start, count int) ([]Card, error) {
 statement := fmt.Sprintf("SELECT id_card, token, bin FROM banwirecard LIMIT %d OFFSET %d", count, start)
    rows, err := db.Query(statement)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    cards := []Card{}
    for rows.Next() {
        var u Card
        if err := rows.Scan(&u.ID, &u.Token, &u.Bin); err != nil {
            return nil, err
        }
        cards = append(cards, u)
    }
    return cards, nil
}