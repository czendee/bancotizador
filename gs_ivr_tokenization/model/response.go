package model


///////////////////response for payment
type ResponsePayment struct {    
	StatusMessage string       `json:"status_message"`
	Status       string       `json:"status"`
	SucessData   ExitoData     `json:"sucess_data"`
}

type ExitoData struct {
	Token  string   `json:"card_token"`
    PaymentReference  string   `json:"paymentreference"`
	Authcode  string   `json:"authcode"`
	Idtransaction  string   `json:"idtransaction"`
	Marca  string   `json:"card_brand"`
	Bin  string   `json:"card_bin"`
    LastDigits  string   `json:"card_last"`
	Type  string   `json:"type_card"`
}

type ExitoDataTokenized struct {
	Token  string   `json:"token"`
    Type  string   `json:"type"`
	Category  string   `json:"category"`
}
