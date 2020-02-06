package main

import (
	"fmt"
	"log"
	"strings"
//	"banwire/services/gs_ivr_tokenization/db"
	modelito "banwire/services/gs_ivr_tokenization/model"
	 "database/sql"
	 _ "github.com/lib/pq"   //use go get github.com/lib/pq
	miu "banwire/services/gs_ivr_tokenization/util"
)

/*
    const (
        DB_USER     = "lerepagr"
        DB_PASSWORD = "Ag8q2utgSsVy2tyR7_M9cNYbzsqSvwma"
        DB_NAME     = "lerepagr"
        DB_SERVER     = "stampy.db.elephantsql.com" //"54.163.207.112"
        DB_PORT      = 5432
    )
 */


///////////////// ///////////////////////////////////////version 3.2


///////////////// ///////////////////////////////////////version 3.2


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
//						DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
						Config_DB_server,Config_DB_name, Config_DB_user, Config_DB_pass, Config_DB_port)
				
				

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



   func logicGeneratetokenizedDBV2(requestData  modelito.RequestTokenized, dataObtained modelito.ExitoDataTokenized ,errorGeneral string) ( modelito.Card,string) {
	////////////////////////////////////////////////process db steps
   //START    
		var miCard modelito.Card
				//  START insert record in Card
				    var errdb error
				    var db *sql.DB
				    // Create connection string
					connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
//						DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
                       Config_DB_server,Config_DB_name, Config_DB_user, Config_DB_pass, Config_DB_port)
				
				

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
						         log.Print("Ping ok!\n")
                                    //verifica si ya existe ese tiken con algun otro cliente
                                    //START
//	                                 var miCard modelito.Card//to return the bin, last, brand, type_card GetCardByToken
							          log.Print(" verificar si ya existe ese token en tabla cards 01!\n")
							         miCard.Token = dataObtained.Token //from the webservice cr.banwire.com method ADD
//							         errCard:= miCard.GetCardByToken(db)
							         errCard:= miCard.GetCardByTokenAndCust(db,miCustomer.ID)							         
							         
							         
						          	log.Print(" verificar si ya existe ese token en tabla cards  para el mismo cliente 02!\n")
									if errCard != nil {
										 if strings.Contains(errCard.Error(),"no rows in result set") {
                                          //no existe, entocnes procede a insertarlo
                                          log.Print(" TOKEN does not exist for the same customer"+errCard.Error())
											//no existe ese token para algun customer reference, proceder a insertar en cards table
											//START
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
												 miCard.Phone =requestData.Phone  //10 digits passed as params
											         errUpdate:=miCard.CreateCard(db)
											          log.Print("regresa func  updateCard ok!\n")
													if errUpdate != nil {
													  log.Print("Error: :"+ errUpdate.Error())
														errorGeneral =errUpdate.Error()
													}
													
												}//end else dataos ok para  card
		                                       //END                                          
                                          //end if strings.contains 
										 }else{
										 	//error de la DB
											  log.Print("Error: Checking token-customer:customer."+errCard.Error() )
											  errorGeneral ="Error: Checking token-customer:TOKEN already exists for this customer."+errCard.Error() 
										 }
										 
									} else{
										
										log.Print(" ya existe table card  token:!\n"+miCard.Token)
										log.Print(" ya existe table card  bin:!\n"+miCard.Bin)
										log.Print(" ya existe table card customer:!\n"+miCard.Customer)
										log.Print(" ya existe table card phone:!\n"+miCard.Phone)
/*									         miCard.Token 
									         miCard.Last 
									         miCard.Bin 
									         miCard.Valid 
									         miCard.Score 
									         miCard.Customer
									         miCard.Brand 
									         miCard.Type 
*/									         
									  log.Print("Error: Checking token-customer:TOKEN already exists for this customer.")
									  errorGeneral ="Error: Checking token-customer:TOKEN already exists for this customer."

								    }
	
                                    //END
                                    


								
								}//end else, no error del select					         


				
					    } //end else de no error ping
				
				
					}//end if no error db
				    
				//  END fetchFromDB

   	  return  miCard, errorGeneral
   }


/////////////////////////////////v4
/////////////////////////////////v4




   func logicDBCheckNumberOfPaymentsToday(requestData modelito.RequestPayment, errorGeneral string) ([]modelito.Payment,string) {
	////////////////////////////////////////////////obtain parms in JSON
   //START    
var resultPayments []modelito.Payment
var errPayments error

				//  START fetchFromDB
				    var errdb error
				    var db *sql.DB
				    // Create connection string
					connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
//						DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
                        Config_DB_server,Config_DB_name, Config_DB_user, Config_DB_pass, Config_DB_port)				
				

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
				
					         resultPayments,errPayments =modelito.GetTodayPaymentsByTokenCard(db,requestData.Token)
					         					         log.Print("regresa func  GetTodayPaymentsByTokenCard ok!\n")
							if errPayments != nil {
							  log.Print("Error: :"+ errPayments.Error())
							  errorGeneral=errPayments.Error()
							}
							var cuantos int
							cuantos = 0
				         	for _, d := range resultPayments {
				         		log.Print("el registor trae:"+d.Token+" "+d.Amount)
							    cuantos =cuantos +1
			         		}
							if cuantos == 0 {
							  log.Print("DB: records not found")
							  
							}else if cuantos >=3{
							  log.Print("DB: Max daily payments for this Credit Card exceeded")
							  errorGeneral="Not more payments can be processed today for this Credit Card. Max number of payments reached"
                                
                            }	

					    }
				
				
					}
				    
				//  END fetchFromDB
   
   //END
   	  return  resultPayments, errorGeneral
   }


   func logicProcesspaymentDBV4(requestData modelito.RequestPayment, errorGeneral string) (modelito.RequestPayment,modelito.Card,string) {
	////////////////////////////////////////////////process db steps
   //START  
           var miCard modelito.Card//to return the bin, last, brand, type_card
           var miPayment modelito.Payment//to insert the payment
		    //  START 
			    var errdb error
			    var db *sql.DB
			    // Create connection string
				connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
//				DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
                Config_DB_server,Config_DB_name, Config_DB_user, Config_DB_pass, Config_DB_port)		
		

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
				         var miCustomer modelito.Customer
				         miCustomer.Reference = requestData.Clientreference 
				         errCustomer:= miCustomer.GetCustomerByReference01(db)
				         //in miCustomer.ID is the value of the id_customer 
						if errCustomer != nil {
							//the customer does not exist to score this payment
							
						  log.Print("Error: Customer does not Exists, payment done, buit score not updated: "+ errCustomer.Error())
						  errorGeneral ="Error: Customer does not Exists. Payment applied, but card score not increased: "+ errCustomer.Error()
                           
						} else{
							//the customer exists
					         log.Print("the customer exists, ID interno es "+miCustomer.ID)
					         miCard.Token =requestData.Token
					         errUpdate:=miCard.IncreaseScoreCardAndCust(db,miCustomer.ID )
					         log.Print("regresa func  IncreaseScoreCard ok!\n")
							 if errUpdate != nil {
								  log.Print("Error: increasing the score for this card:"+ errUpdate.Error())
							      errorGeneral =errUpdate.Error()
 							 }else{
                                    //the increase was done, now try record the payment for rule (3max payments for tcd a day)
                                    log.Print("About record Payment info in DB, the customer exists, ID interno es "+miCustomer.ID)
                                    miPayment.Token =requestData.Token
                                    miPayment.Amount =requestData.Amount
				    miPayment.Msi =requestData.Msi
                                    errInsertPay:=miPayment.CreatePayment(db )
                                    log.Print("regresa func  CreatePayment ok!\n")
                                    if errInsertPay != nil {
                                        log.Print("Error: Recording the payment info in the DB:"+ errInsertPay.Error())
                                        errorGeneral =errInsertPay.Error()
                                    }else{
                                        //increase ok and the payment detail recorded  ok
                                        log.Print(" se ejecuta  select table card to get bin, last, brand. type  01!\n")
                                            miCard.Token = requestData.Token
                                            errCard:= miCard.GetCardByToken(db)
                                        log.Print(" se ejecuta select table card to get bin, last, brand. type  02!\n")
                                            if errCard != nil {
                                            log.Print("Error: after payment was applied and score increased,There was a problem getting the customer:"+ errCard.Error())
                                            errorGeneral ="Error: after payment was applied and score increased,There was a problem getting the customer:"+ errCard.Error()
                                            
                                            } else{
                                              log.Print(" select table card to get token:!\n"+miCard.Token)
                                              log.Print(" select table card to get bin:!\n"+miCard.Bin)
                                              log.Print(" select table card to get last:!\n"+miCard.Last)
                                            }

                                    }                                  
	
							 }//end else de increase

                        }//end else de customer does exists

			
				    }
			
			
				}
		    
		//  END updateCardScoreDB
   
   	  return  requestData, miCard, errorGeneral
   }

   func logicDBRemoveCardIfNotPreviousPayment(requestData modelito.RequestPayment, errorGeneral string) (string,string) {
	//////////////////////////////////////////////
   //START    

 var miCard modelito.Card//to return the bin, last, brand, type_card
var errPayments error
var resultMssg string

resultMssg=""
				//  START fetchFromDB
				    var errdb error
				    var db *sql.DB
				    // Create connection string
					connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
//						DB_SERVER,DB_NAME, DB_USER, DB_PASSWORD, DB_PORT)
                        Config_DB_server,Config_DB_name, Config_DB_user, Config_DB_pass, Config_DB_port)				
				

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
                            //start else -get previous paymnets to check if card is deleted or not
					         log.Print("Ping ok!\n")
				
					         resultMssg,errPayments =modelito.GetAllPaymentsByTokenCard(db,requestData.Token)
					         					         log.Print("regresa func  GetTodayPaymentsByTokenCard ok!\n")
							if errPayments != nil {
							  log.Print("Error: :"+ errPayments.Error())
							  errorGeneral="Error: NO ABLE TO CHECK PREVIOUS PAYMENTS"+ errPayments.Error()
							}
							if resultMssg != "" {
                                if resultMssg == "EXISTEN PREVIOS"{
                                    //no se puede borrar
                                    resultMssg ="Card not removed, as there are previous payments. OK"
                                }else{
                                    //se debe borrar la card
                                    //the increase was done, now try record the payment for rule (3max payments for tcd a day)
                                    log.Print("Delete card, as there was a problem with the payment and there are not previous payments ")
                                    miCard.Token =requestData.Token

                                    errDeleteCard:=miCard.DeleteCard(db )
                                    log.Print("regresa func  DeleteCard ok!\n")
                                    if errDeleteCard != nil {
                                        log.Print("Error: Deleting card info in the DB:"+ errDeleteCard.Error())
                                        errorGeneral ="Error: Deleting card info in the DB:"+ errDeleteCard.Error()
                                    }else{
                                        //card was removed succesfully,  ok
                                        resultMssg ="Card removed OK"
                                    }                                    
                                }//end else Delete Card
                                

							}else{//End if, hay mensaje 
                                //resultMssg is nil, as ther was an error
                                resultMssg=""
                            }

					      }//end else -get previous paymnets to check if card is deleted or not
				
				
					}
				    
				//  END db connDB
   
   //END
   	  return  resultMssg, errorGeneral
   }
