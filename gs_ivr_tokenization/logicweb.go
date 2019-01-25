package main

import (
//	"fmt"
	"log"

//	"banwire/services/gs_ivr_tokenization/net"
	modelito "banwire/services/gs_ivr_tokenization/model"
	miu "banwire/services/gs_ivr_tokenization/util"
//	"time"
//	"encoding/json"
//	 "database/sql"
	 _ "github.com/lib/pq"   //use go get github.com/lib/pq
	 
    "strings"
//    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
)


///////////////////////////////////////business logic functions

//////////////////////////////////web


   func logicProcesspaymentWeb(requestData modelito.RequestPayment, errorGeneral string) (modelito.ExitoData,string) {
   	
   	     var resultadoPayment modelito.ExitoData  //response.go

	    valoresbanwire := url.Values{
		"method": {"payment"},
//		"user": {"pruebasbw"},  //this value was mentioned by Charly dec 05,2018
		"user": {Config_WS_crbanwire_pass},  //this value needs to be configurable to Move to production. 22 Jan 2019        


		"reference": {requestData.Paymentreference}, 
		"token": {requestData.Token},         
		"amount": {requestData.Amount},          
	    "cvv": {requestData.Cvv},  
	}
log.Println("web api el cvv"+requestData.Cvv)

	    response,err := http.PostForm("https://cr.banwire.com/?action=card",
	valoresbanwire)
	
	
	    if err != nil {
	        log.Printf("The HTTP request failed with error %s\n", err)
	    } else {
	        data, _ := ioutil.ReadAll(response.Body)
	        log.Println(string(data))
	        if strings.Contains(string(data), "error") {
	         	errorGeneral ="Error returned: "+string(data)
	        }else{
		        	var str WebResponsePayment
					_ = json.Unmarshal(data, &str)
					
					if str.Description == "Payment on Demand" {
					    // Do Stuff
					    resultadoPayment.Token  =requestData.Token
					    resultadoPayment.PaymentReference = str.Reference
					    resultadoPayment.Authcode =str.Authcode
					    resultadoPayment.Idtransaction =str.Idtransaction

					}else{
	         			errorGeneral ="Error in JSON returned/internal webservice"

					}
	        	
	        }
	    }
	    log.Println("web api Terminating the application...")
	    
   	  return  resultadoPayment, errorGeneral
   }
   
   
   func logicGeneratetokenizedWeb(requestData modelito.RequestTokenized, errorGeneral string) (modelito.ExitoDataTokenized,string) {
		////////////////////////////////////////////////process db steps
	   //START    

      var resultadoService modelito.ExitoDataTokenized

     month,year, errorGeneral := miu.ConvertMMYYintoMonthYear(requestData.Exp)   //utils.go

      
	   if errorGeneral==""{
		    valoresbanwire := url.Values{
			"method": {"add"},
//		"user": {"pruebasbw"},  //this value was mentioned by Charly dec 05,2018
		"user": {Config_WS_crbanwire_pass},  //this value needs to be configurable to Move to production. 22 Jan 2019        
			"email": {"generalseguros@genearlseguros.com"},
			"number": {requestData.Card},  //
			"exp_month": {month},              //requestData.Exp  solo mes
			"exp_year": {year},             //requestData.Exp  solo año u de 4 digitos
//05Dec2018	"cvv": {requestData.Cvv},                   //iissue necesita el cvv
			"cvv": {"000"},                   //Charly V indicated to send 000
			"name": {"generalseguros"}, 
			"address": {"generalseguros"},
			"postal_code": {"06000"},
			}	

		    response,err := http.PostForm("https://cr.banwire.com/?action=card&exists=1",
			valoresbanwire)
			
		    if err != nil {
		        log.Printf("The HTTP request failed with error %s\n", err)
		         errorGeneral ="The HTTP request failed with error "+err.Error()
		    } else {
		        data, _ := ioutil.ReadAll(response.Body)
		        log.Println("webservice response:" +string(data))
		        if strings.Contains(string(data), "error") {
		         	errorGeneral ="Error returned: "+string(data)
		        }else{
		        	log.Println("logicGeneratetokenizedWeb 02")

		        	var str WebResponseAdd
                      
					_ = json.Unmarshal(data, &str)
		        	log.Println("logicGeneratetokenizedWeb 03")
   					log.Println("logicGeneratetokenizedWeb token:"+str.Token)
   					log.Println("logicGeneratetokenizedWeb task:"+str.Task)
   					log.Println("logicGeneratetokenizedWeb type:"+str.Card.Type+"::")
   					log.Println("logicGeneratetokenizedWeb cat:"+str.Card.Category+"::")
				    resultadoService.Token  =str.Token
				    resultadoService.Type  =str.Card.Type
				    resultadoService.Category  =str.Card.Category
///revisar con charly y toño estas reglas
/// cuando si insertar y cuando no
				    
					if str.Task == "add" {
		        		log.Println("logicGeneratetokenizedWeb 04")
						if str.Result {
		        			log.Println("logicGeneratetokenizedWeb 05")
							if str.Exists  {
								//error
								 // the card already exists
		         				//aun asi, debe de registrarse en DB,
		         				//errorGeneral ="Error : this card already exists and was tokenized previously"
							    resultadoService.Token  =str.Token
							    resultadoService.Type  =str.Card.Type
							    resultadoService.Category  =str.Card.Category
   					        	log.Println("logicGeneratetokenizedWeb 07:"+resultadoService.Token)
   					        	log.Println("logicGeneratetokenizedWeb 07:"+resultadoService.Type)

							}else{ // this is the flow we care
   					        	log.Println("logicGeneratetokenizedWeb 06")
					    // Do Stuff
							    resultadoService.Token  =str.Token
							    resultadoService.Type  =str.Card.Type
							    resultadoService.Category  =str.Card.Category
   					        	log.Println("logicGeneratetokenizedWeb 07:"+resultadoService.Token)
   					        	log.Println("logicGeneratetokenizedWeb 07:"+resultadoService.Type)
							}
							
						}else{
							// the card tokenization fail
		         			errorGeneral ="Error : the result of the tokenization is FALSE"
						}


					}else{
						// the task for the card needs to be add
		         		errorGeneral ="Error : intenal configuration task needs to be ADD"
					}
					
		        }
		    }
	   	
	   }
     
	    log.Println("web api Terminating the application...")
   	  return  resultadoService, errorGeneral
   }
   
   
 type WebResponseAdd struct {
//    ID string `json:"id"`
    Token string `json:"token"`
    Task string `json:"task"`
    Result bool `json:"result"`
    Exists bool `json:"exists"`
    Card card `json:"card"`
   
}

 type card struct {
    Type string `json:"type"`
    Category string `json:"category"`    
}

 type WebResponsePayment struct {
    ID string `json:"id"`
    Idtransaction string `json:"id_transaction"`
    Authcode string `json:"auth_code"`
    Reference string `json:"reference"`
    Description string `json:"description"`
    Amount string `json:"amount"`
}