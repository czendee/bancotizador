package main

import (
	"net/http"
	"log"
	"banwire/services/gs_ivr_tokenization/db"
//	"banwire/services/gs_ivr_tokenization/net"
	modelito "banwire/services/gs_ivr_tokenization/model"
//	"time"
//	"encoding/json"
	 _ "github.com/lib/pq"   //use go get github.com/lib/pq

)


    
func ProcessGettokenizedcards(w http.ResponseWriter,  requestData modelito.RequestTokenizedCards) (string,string) {

    var errorGeneral string
    var errorGeneralNbr string
    	var result string
   var valoresParaResponder  []modelito.Card

    errorGeneral=""

    if errorGeneral=="" {


		log.Print("CZ   STEP Validate Parms")

		/// START
	        result ="OK realizarpago"+requestData.Cardreference
		    log.Print("CZ    handler Listening fetch:"+result)
		    
		     log.Print("CZ   STEP Validate paramters request")
		    errorGeneral= validaReqFetchCards(requestData)
		
		
		/// END

    }				    
		        
    if errorGeneral!="" && errorGeneralNbr=="" {
    	//prepare response with error 300
    	log.Print("CZ    Prepare Response with 300. Missing parameter:"+errorGeneral)
    	errorGeneral="ERROR:300 -Missing parameter"	+errorGeneral
    	errorGeneralNbr="300"
    }

	////////////////////////////////////////////////DB	
	//	    resultado,errfetchDB:= fetchFromDB ()
	if errorGeneral==""{//continue next step

       	    	log.Print("CZ   STEP Consume DB")
         valoresParaResponder,errorGeneral =logicDBGettokenizedcardsV2(requestData, errorGeneral) 


    }				    
    if errorGeneral!="" && errorGeneralNbr==""{
    	//prepare response with error 310
    	log.Print("CZ    Prepare Response with 310. Error obtaining cards:"+errorGeneral)
    	errorGeneral="ERROR:310 -  Error obtaining cards -"	+errorGeneral
	    errorGeneralNbr="310"
    }

		        
	//response
    log.Print("CZ    handler DB Listening test gettokenizedcards  2")					

	//////////    format the response
    if errorGeneral==""{//continue next step
		log.Print("CZ   STEP Validate Parms")
			/// START
		fieldDataBytesJson,err := getJsonResponseV2(valoresParaResponder)
		
		log.Print("CZ    handler Listening test gettokenizedcards  3")	
		
		result ="OK gettokenizedcards"+requestData.Cardreference+"resultado"
		//////////    write the response
		w.Header().Set("Content-Type", "application/json")
		 w.Write(fieldDataBytesJson)
		 
		 log.Print("CZ    handler Listening test gettokenizedcards  4"+"<html><body>"+ result+"</body></html>")
			         
        if err!=nil{
        	log.Print("Eror en generando response")
            errorGeneral= err.Error()
        }		
		
		/// END

    }				    
		 
    if errorGeneral!="" && errorGeneralNbr==""{//continue next step
    	log.Print("CZ   prepare the JSON response for ERROR")

	    //  START 
	    errorGeneral="ERROR:330 -Error preparing the response"	+errorGeneral
	    errorGeneralNbr="330"
	    //  END
    }

     return errorGeneral, errorGeneralNbr
}



// Generatetokenized for receive and handle the request from client
func ProcessGeneratetokenized(w http.ResponseWriter, requestData modelito.RequestTokenized) (string,string) {
	defer func() {
		db.Connection.Close(nil)
	}()
	  var result string

     var errorGeneral string
     var errorGeneralNbr string
     
     var resultCardTokenized modelito.Card
     
     var obtainedDataWebservice modelito.ExitoDataTokenized
     
    errorGeneral=""


	////////////////////////////////////////////////validate parms
	/// START
    if errorGeneral==""{//continue next step
		        result ="OK realizarpago"+requestData.Clientreference+"    :    " +requestData.Paymentreference+"    :    " +requestData.Card+"    :    " +requestData.Exp+"    :    " +requestData.Cvv
    		    log.Print("CZ    handler Listening test handleGeneratetokenized:"+result)
		     log.Print("CZ   STEP Validate paramters request")
		    errorGeneral= validaReqGenerateTokenized(requestData)	
		/// END

	}	
		              
    if errorGeneral!="" && errorGeneralNbr=="" {
    	//prepare response with error 800
    	log.Print("CZ    Prepare Response with 200. Missing parameter:"+errorGeneral)
    	errorGeneral="ERROR :200 -Missing parameter "	+errorGeneral
		errorGeneralNbr="200"
    }




	////////////////////////////////////////////////consume internal websrvice banwire
	//////////////////            tokenization 

    if errorGeneral==""{//continue next step
				/// START
				obtainedDataWebservice, errorGeneral =logicGeneratetokenizedWeb(requestData, errorGeneral)
				
				/// END
	}	

    if errorGeneral!="" && errorGeneralNbr==""{
    	//prepare response with error 210
    	log.Print("CZ    Prepare Response with 210. Error processing payment:"+errorGeneral)
    	errorGeneral="ERROR:210 -Error processing payment:"	+errorGeneral
		errorGeneralNbr="210"
    }

				
				
	////////////////////////////////////////////////DB	
	//	    insert new record in Card , if customer doesn't exist, insert a new one?
	//  Update if exist, if not insert in Customer

    if errorGeneral==""{//continue next stepjhlkjg 
        	log.Print("CZ   el  token:"+obtainedDataWebservice.Token)
    				resultCardTokenized, errorGeneral =logicGeneratetokenizedDBV2(requestData,obtainedDataWebservice , errorGeneral)
    						
	}					

    if errorGeneral!="" && errorGeneralNbr==""{
    	//prepare response with error 220
    	log.Print("CZ    Prepare Response with 220. Error generating token:"+errorGeneral)
    	errorGeneral="ERROR:220 -Error generating token:"	+errorGeneral
		errorGeneralNbr="220"
    }

	//response
    if errorGeneral==""{//continue next step
		log.Print("CZ   STEP Post the response JSON ready")
		
			/// START
		fieldDataBytesJsonTokenize,err := getJsonResponseTokenizeV2(resultCardTokenized)
			
		log.Print("CZ    handler Listening test realizarpago  3")	
	    
	    w.Header().Set("Content-Type", "application/json")
	    w.Write(fieldDataBytesJsonTokenize)
		log.Print("CZ    handler Listening test handleGeneratetokenized  4"+"<html><body>"+ result+"</body></html>")		         		         
        if err!=nil{
        	log.Print("Eror en generando response")
	        errorGeneral= err.Error()
        }
				
		/// END
	}	

    if errorGeneral!="" && errorGeneralNbr==""{
    	//prepare response with error 230
    	log.Print("CZ    Prepare Response with 230. Error generating Response Tokenized:"+errorGeneral)
    	errorGeneral="ERROR:230 -Error generating Response Tokenized:"	+errorGeneral
		errorGeneralNbr="230"
    }
    
	 log.Print("CZ  ends func tokenized")
	 
	return errorGeneral, errorGeneralNbr
}





func GetCardType(number string) string {
	return "VISA"
}


/////////////////////////v4
/////////////////////////v4

// v4Processpayment  receive and handle the request from client, access DB
func v4ProcessProcessPayment(w http.ResponseWriter, requestData modelito.RequestPayment) (string,string){
	defer func() {
		db.Connection.Close(nil)
	}()
    var result string
    var errorGeneral string
    var	errorGeneralNbr string

    var resultadoPayment modelito.ExitoData
    errorGeneral=""

	////////////////////////////////////////////////validate parms
	/// START
    
    if errorGeneral==""{//continue next step
    	log.Print("CZ   STEP Validate Parms")

		/// START
	        result ="OK realizarpago"+requestData.Clientreference+"    :    " +requestData.Paymentreference+"    :    " +requestData.Token+"    :    " +requestData.Cvv+"    :    " +requestData.Amount+"    :    "
		    log.Print("CZ    handler Listening test realizarpago:"+result)
		    
		     log.Print("CZ   STEP Validate paramters request")
		    errorGeneral= validaReqProcessPayment(requestData)
		
		
		/// END

    }				    
		        
    if errorGeneral!="" {
    	//prepare response with error 100
    	log.Print("CZ    Prepare Response with 100. Missing parameter:"+errorGeneral)
    	errorGeneral="ERROR:100 - Missing parameter"	+errorGeneral
		errorGeneralNbr="100"
    }
//////////////////////////////////////////DB verify if less payments for the same card
//////////////////////////////////////////in the same day

	////////////////////////////////////////////////DB	
	//	    resultado,errfetchDB:= fetchFromDB ()
 var valoresParaResponder  []modelito.Payment
	if errorGeneral==""{//continue next step

     	 log.Print("CZ   STEP Consume DB to check if more payments cvan be done today for this card")
         valoresParaResponder,errorGeneral =logicDBCheckNumberOfPaymentsToday(requestData, errorGeneral) 

    }				    
    if errorGeneral!="" && errorGeneralNbr==""{
    	//prepare response with error 105
    	log.Print("CZ    Prepare Response with 105. Error Max payments today for this card exceeded:"+errorGeneral)
    	errorGeneral="ERROR 105 -  Error Max payments today for this card exceeded -"	+errorGeneral
	    errorGeneralNbr="105"
    }
     if valoresParaResponder == nil{
         
     }
		        
	//response
    log.Print("CZ    handler DB Listening test gettokenizedcards  2")					

	////////////////////////////////////////////////consume internal websrvice banwire
	//////////////////            process payment
	    if errorGeneral==""{//continue next step
	    	log.Print("CZ   STEP Consume internal websrvice banwire")

			/// START
			
			resultadoPayment, errorGeneral= logicProcesspaymentWeb(requestData , errorGeneral )  
			/// END

	    }				    
	    if errorGeneral!="" && errorGeneralNbr==""{
	    	//prepare response with error 110
	    	log.Print("CZ    Prepare Response with 110. Error processing payment:"+errorGeneral)
	    	errorGeneral="ERROR:110 - Error processing payment"	+errorGeneral
			errorGeneralNbr="110"
	    }
								

	////////////////////////////////////////////////DB	
	//      update the score field: increase by 1
	//      for this card
	//	    
	var  dataObtainedCard  modelito.Card
	    if errorGeneral==""{//continue next step
	    	log.Print("CZ   STEP  update the score field: increase by 1")
			requestData, dataObtainedCard, errorGeneral= logicProcesspaymentDBV4(requestData , errorGeneral )  

									log.Print(" medio token:!\n"+dataObtainedCard.Token)
									log.Print(" medio bin:!\n"+dataObtainedCard.Bin)
									log.Print(" medio last:!\n"+dataObtainedCard.Last)
		    resultadoPayment.Marca = dataObtainedCard.Brand
		    resultadoPayment.Bin = dataObtainedCard.Bin
		    resultadoPayment.LastDigits= dataObtainedCard.Last
		    resultadoPayment.Type = dataObtainedCard.Type
		    
	    }				    

	    if errorGeneral!="" && errorGeneralNbr==""{
	    	//prepare response with error 120
	    	log.Print("CZ    Prepare Response with 120. Error recording results in DB:"+errorGeneral)
	    	errorGeneral="ERROR: 120 - Error recording results in DB"	+errorGeneral
			errorGeneralNbr="120"
	    }

    		    
	//response
	////////////////////////////////////////////////http response	
	//      prepare the JSON response
	//	    
	    if errorGeneral==""{//continue next step
	    	log.Print("CZ   STEP  prepare the JSON response for SUCCESS")

		    //  START 

		    fieldDataBytesJsonPayment,err := getJsonResponsePaymentV2(resultadoPayment)					
		        w.Header().Set("Content-Type", "application/json")
		        w.Write(fieldDataBytesJsonPayment)
				log.Print("CZ    handler Listening test handleProcesspayment  4"+"<html><body>"+ result+"</body></html>")		         
                if err!=nil{
                	log.Print("Eror en generando response")
                    errorGeneral= err.Error()
                }
		    //  END
        }

	    if errorGeneral!="" && errorGeneralNbr=="" {//continue next step
	    	log.Print("CZ   prepare the JSON response for ERROR")

		    //  START 
		    errorGeneral="ERROR:130 -Error preparing the response"	+errorGeneral
			errorGeneralNbr="130"
		    //  END
        }
 log.Print("CZ  END   handler Listening DB  realizarpago  2")	
     return errorGeneral, errorGeneralNbr
}
