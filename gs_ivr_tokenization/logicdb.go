package main

import (
	"fmt"
	"log"
//	"banwire/services/gs_ivr_tokenization/db"
	modelito "banwire/services/gs_ivr_tokenization/model"
	 "database/sql"
	 _ "github.com/lib/pq"   //use go get github.com/lib/pq
	miu "banwire/services/gs_ivr_tokenization/util"
)


    const (
        DB_USER     = "banwire"
        DB_PASSWORD = "banwire"
        DB_NAME     = "banwire"
        DB_SERVER     = "bin_banwire_service_gs_ivr_postgres" //"54.163.207.112"
        DB_PORT      = 5432
    )
 


///////////////// ///////////////////////////////////////version 2


///////////////// ///////////////////////////////////////version 2


   func logicDBGettokenizedcardsV2(requestData modelito.RequestTokenizedCards, errorGeneral string) ([]modelito.Card,string) {
	////////////////////////////////////////////////obtain parms in JSON
   //START    
var resultCards []modelito.Card
var errCards error

				//  START fetchFromDB
				    var errdb error
				    var db *sql.DB
				    // Create connection string
					connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
						DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
				
				

					 // Create connection pool
					db, errdb = sql.Open("postgres", connString)
					if errdb != nil {
						log.Print("Error creating connection pool: " + errdb.Error())
						errorGeneral=errdb.Error()
					}
					// Close the database connection pool after program executes
					 defer db.Close()
					if errdb == nil {
						log.Print("Connected!\n")
				
					
						errPing := db.Ping()
						if errPing != nil {
						  log.Print("Error: Could not establish a connection with the database:"+ errPing.Error())
							  errorGeneral=errPing.Error()
						}else{
					         log.Print("Ping ok!\n")
//					         var misCards modelito.Card
					         
					         resultCards,errCards =modelito.GetCardsByCustomer(db,requestData.Cardreference)
					         					         log.Print("regresa func  getCardsByCustomer ok!\n")
							if errCards != nil {
							  log.Print("Error: :"+ errCards.Error())
							  errorGeneral=errCards.Error()
							}
							var cuantos int
							cuantos = 0
				         	for _, d := range resultCards {
				         		log.Print("el registor trae:"+d.Token+" "+d.Bin)
							    cuantos =1
			         		}
							if cuantos == 0 {
							  log.Print("DB: records not found")
							  errorGeneral="Not cards found for the customer reference received"
							}		

					    }
				
				
					}
				    
				//  END fetchFromDB
   
   //END
   	  return  resultCards, errorGeneral
   }

   func logicProcesspaymentDBV2(requestData modelito.RequestPayment, errorGeneral string) (modelito.RequestPayment,modelito.Card,string) {
	////////////////////////////////////////////////process db steps
   //START  
           var miCard modelito.Card//to return the bin, last, brand, type_card
		    //  START 
			    var errdb error
			    var db *sql.DB
			    // Create connection string
				connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
				DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
		
		

			 // Create connection pool
				db, errdb = sql.Open("postgres", connString)
				if errdb != nil {
					log.Print("Error creating connection pool: " + errdb.Error())
				}
				// Close the database connection pool after program executes
				 defer db.Close()
				if errdb == nil {
					log.Print("Connected!\n")
			
				
					errPing:= db.Ping()
					if errPing != nil {
					  log.Print("Error: Could not establish a connection with the database:"+ errPing.Error())
					  errorGeneral =errPing.Error()
					}else{
				         log.Print("Ping ok!\n")

				         miCard.Token =requestData.Token
				         errUpdate:=miCard.IncreaseScoreCard(db)
				          log.Print("regresa func  IncreaseScoreCard ok!\n")
						if errUpdate != nil {
						  log.Print("Error: :"+ errUpdate.Error())
					      errorGeneral =errUpdate.Error()
						}else{
					          log.Print(" se ejecuta  select table card to get bin, last, brand. type  01!\n")
						         miCard.Token = requestData.Token
						         errCard:= miCard.GetCardByToken(db)
				          	log.Print(" se ejecuta select table card to get bin, last, brand. type  02!\n")
								if errCard != nil {
								  log.Print("Error: get customer:"+ errCard.Error())
								  errorGeneral =errCard.Error()
	                               
								} else{
									log.Print(" select table card to get token:!\n"+miCard.Token)
									log.Print(" select table card to get bin:!\n"+miCard.Bin)
									log.Print(" select table card to get last:!\n"+miCard.Last)
							    }

						}
			
				    }
			
			
				}
		    
		//  END updateCardScoreDB
   
   	  return  requestData, miCard, errorGeneral
   }


   func logicGeneratetokenizedDBV2(requestData  modelito.RequestTokenized, dataObtained modelito.ExitoDataTokenized ,errorGeneral string) ( modelito.Card,string) {
	////////////////////////////////////////////////process db steps
   //START    
		var miCard modelito.Card
				//  START insert record in Card
				    var errdb error
				    var db *sql.DB
				    // Create connection string
					connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
						DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
				
				

					 // Create connection pool
					db, errdb = sql.Open("postgres", connString)
					if errdb != nil {
						log.Print("Error creating connection pool: " + errdb.Error())
						errorGeneral =errdb.Error()
					}
					// Close the database connection pool after program executes
					 defer db.Close()
					if errdb == nil {
						log.Print("Connected!\n")
				
					
						errPing:= db.Ping()
						if errPing != nil {
						  log.Print("Error: Could not establish a connection with the database:"+ errPing.Error())
						  errorGeneral= errPing.Error()
						}else{
						         log.Print("Ping ok!\n")
						         var miCustomer modelito.Customer
						         miCustomer.Reference = requestData.Clientreference 
						         errCustomer:= miCustomer.GetCustomerByReference01(db)
						         //in miCustomer.ID is the value of the id_customer 
								if errCustomer != nil {
								  log.Print("Error: get customer:"+ errCustomer.Error())
								  errorGeneral =errCustomer.Error()
	                               
								} else{
							             log.Print("Listo para insertar card!\n")
								         milast,errLast :=miu.ObtainLast4fromCard (requestData.Card) //utils.go
								         mibin,errBin :=miu.ObtainBINfromCard (requestData.Card) //utils.go
										if(errLast!=""){
											errorGeneral =errLast
							                log.Print("error obatining the last 4!\n")
										}else if(errBin!=""){
											errorGeneral =errBin
							                log.Print("error obatining the BIN!\n")
										}else{
							                log.Print(" todo ok para insertar!\n")
									         miCard.ID ="888"   //current value +1  o un random
							                log.Print(" todo ok para insertar, el parametro de token es !\n"+dataObtained.Token +"   : "+dataObtained.Type)
									         miCard.Token =dataObtained.Token// value returned by the internal webservice 
									         miCard.Last =milast//ulitmos 4 digitos de card
									         miCard.Bin =mibin //6 basic digits in a card
									         miCard.Valid =requestData.Exp  //4 digits passed as params
									         miCard.Score ="0"
									         miCard.Customer =miCustomer.ID	
									         miCard.Brand = miu.GetCardType(requestData.Card)
									         miCard.Type = dataObtained.Type
									         errUpdate:=miCard.CreateCard(db)
									          log.Print("regresa func  updateCard ok!\n")
											if errUpdate != nil {
											  log.Print("Error: :"+ errUpdate.Error())
												errorGeneral =errUpdate.Error()
											}
											
										}//end else dataos ok para  card


								
								}//end else, no error del select					         


				
					    } //end else de no error ping
				
				
					}//end if no error db
				    
				//  END fetchFromDB

   	  return  miCard, errorGeneral
   }
