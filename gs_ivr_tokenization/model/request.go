package model
import (
	"time"
	)


// ViewRequest is the set of parameters for view request
type ViewRequest struct {
	Condition map[string]interface{} `json:"condition"`
	Order     map[string]int         `json:"order"`
	Limit     *int64                 `json:"limit"`
	Offset    *int64                 `json:"offset"`
}




//////////////////http request



//request

type RequestTokenizedCards struct {
	Cardreference string            `json:"cardreference"`
	
}

type RequestPayment struct {
	Clientreference string            `json:"clientreference"`
	Paymentreference  string      `json:"paymentreference"`
	Token  string      `json:"Token"`
	Cvv  string      `json:"Cvv"`
	Amount  string      `json:"Amount"`
}

type RequestTokenized struct {

	Clientreference string            `json:"clientreference"`
	Paymentreference string            `json:"paymentreference"`
	Card  string      `json:"card"`
	Exp  string      `json:"exp"`
	Cvv  string      `json:"Cvv"`
}


type myDataService01 struct {
	Idreference string            `json:"Idreference"`
	Idcomercio  string      `json:"Idcomercio"`
}



type AutoGenerated struct {
	StatusMessage string       `json:"status_message"`
    Status string       `json:"status"`
	Cards       []CardData `json:"card_data"`
}

type CardData struct {
	Date time.Time `json:"card_date"`
	Token  string   `json:"card_token"`
    LastDigits  string   `json:"card_last"`
	Marca  string   `json:"card_brand"`
	Vigencia  string   `json:"card_exp"`
	Bin  string   `json:"card_bin"`
	Score  string   `json:"card_score"`
	Type  string   `json:"type_card"`

}

/////////////////////response for tokenize

type ResponseTokenize struct {
	StatusMessage string       `json:"status_message"`
	Status       string       `json:"status"`
    Card      CardData `json:"card_data"`
}

type ResponseError struct {
	StatusMessage string       `json:"status_message"`
    Status string       `json:"status"`
}