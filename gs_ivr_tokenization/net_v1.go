package main

import (
	"net/http"
	"fmt"
	"log"
	"banwire/services/gs_ivr_tokenization/db"
	"banwire/services/gs_ivr_tokenization/net"
	modelito "banwire/services/gs_ivr_tokenization/model"

    "banwire/services/gs_ivr_tokenization/model"
	"banwire/services/gs_ivr_tokenization/model/pgsql"
    "encoding/json"
    
//	"time"
//	"encoding/json"
//	 "database/sql"
	 _ "github.com/lib/pq"   //use go get github.com/lib/pq
)

    
// init loads the routes for version 1
func init() {
//	var _r = net.GetRouter()
//	var r = _r.PathPrefix("/v1").Subrouter()

    var r = net.GetRouter()
	//route for test
	    log.Print("cz  init net_v1")
	r.Handle("/v3/fetchtokenizedcards", netHandle(handleDBGettokenizedcards, nil)).Methods("GET")   //logicbusiness.go
    r.Handle("/v3/processpayment", netHandle(v4handleDBProcesspayment, nil)).Methods("GET")              //logicbusiness.go 
	r.Handle("/v3/generatetokenized", netHandle(handleDBGeneratetokenized, nil)).Methods("GET")     //logicbusiness.go
	r.Handle("/v3/fetchtokenizedcards", netHandle(handleDBPostGettokenizedcards, nil)).Methods("POST")   //logicbusiness.go
	r.Handle("/v3/processpayment", netHandle(v4handleDBPostProcesspayment, nil)).Methods("POST")           //logicbusiness.go    	    

	r.Handle("/v3/generatetokenized", netHandle(handleDBPostGeneratetokenized, nil)).Methods("POST")     //logicbusiness.go

	r.Handle("/v3/testdbcharly", netHandle(handleTestV1, nil)).Methods("GET")     //logicbusiness.go

}

// handleTest is an example for receive and handle the request from client
func handleTestV1(w http.ResponseWriter, r *http.Request) {
	var response []byte
	var err error

	defer func() {
		db.Connection.Close(r)
		if err != nil {
			err = model.UnmarshalJSONError(err)
			response, _ = json.Marshal(err)
		}
		if response != nil {
			rw := net.ResponseWriterJSON(w)
			rw.Write(response)
		}
	}()

    log.Print("cz  handleTestV1")

	fmt.Fprint(w,"youtochi   antes de  probar la db")
	
	_pg := pgsql.Db{
		db.Connection.GetPgDb(r, false),
	}

	fmt.Fprint(w,"youtochi   establece conn - probar la db")

    //llama una func que acceda la db
    var miCardResult model.Card
    fmt.Fprint(w,"     youtochi   establece conn - 2   ")
    err, miCardResult= _pg.GetCardByTokenAndCust("crd.2Gm9cob47gMEet9hjcMbVS3oJwmf","14")
    


	fmt.Fprint(w,"   youtochi   ver si ok  - probar la db       ")
    log.Println(" 3  youtochi   ver si ok  - probar la db!      \n")
    


	if err != nil {
		var e = model.NewErrorGroup()
		if _e, ok := err.(*pgsql.Error); ok {

			switch _e.Type {
			case pgsql.ErrNotFound:
//				e.Push(model.InvalidValueError("id", error_not_found))
                 log.Print("ERROR carloitos Not Found!\n")
				break
//			case pgsql.ErrNotFoundStatusPayment:
//				e.Push(model.InvalidValueError("id_status_payment", "The status payment id does not exists"))
//				break
			default:
				log.Panic(_e.Type)
				break
			}

		} else {
			log.Panic(err)
		}

		if e.HasError() {
			err = e
			return
		}
	}else{
       log.Println("        youtochi   quesque si  - probar la db!      \n")
            
        log.Print("       youtochi   ver si ok  - ver el dato!\n"+miCardResult.Token)
    }
        
	fmt.Fprint(w,"youtochi   establece conn - probar la db")
    
	rw := net.ResponseWriterJSON(w)
	rw.Write([]byte(`{"readyCarlitos":true}`))
}

   //post
   
   // handleDBGettokenizedcards  receive and handle the request from client, access DB, and web
func handleDBPostGettokenizedcards(w http.ResponseWriter, r *http.Request) {
	defer func() {
		db.Connection.Close(nil)
	}()
    var errorGeneral string
    var errorGeneralNbr string
    
   	var requestData modelito.RequestTokenizedCards


    errorGeneral=""
    requestData, errorGeneral=obtainPostParmsGettokenizedcards(r,errorGeneral) //logicrequest_post.go

	////////////////////////////////////////////////process business rules
	/// START
    if errorGeneral=="" {

		errorGeneral,errorGeneralNbr= ProcessGettokenizedcards(w , requestData)
	}
	/// END
    if errorGeneral!=""{
    	//send error response if any
    	//prepare an error JSON Response, if any
		log.Print("CZ   STEP Get the ERROR response JSON ready")
		
			/// START
		fieldDataBytesJson,err := getJsonResponseError(errorGeneral, errorGeneralNbr)
		//////////    write the response (ERROR)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fieldDataBytesJson)	
		if(err!=nil){
			
		}
	
    } 
					
}




// handleGeneratetokenized for receive and handle the request from client
func handleDBPostGeneratetokenized(w http.ResponseWriter, r *http.Request) {
	defer func() {
		db.Connection.Close(nil)
	}()
     var requestData modelito.RequestTokenized
     var errorGeneral string
     var errorGeneralNbr string
     
    errorGeneral=""


    requestData,errorGeneral =obtainPostParmsGeneratetokenized(r,errorGeneral)   //logicrequest_post.go



	////////////////////////////////////////////////validate parms
	/// START
    
    if errorGeneral=="" {

		errorGeneral,errorGeneralNbr= ProcessGeneratetokenized(w , requestData)
	}

    if errorGeneral!=""{
    	//send error response if any
    	//prepare an error JSON Response, if any
		log.Print("CZ   STEP Get the ERROR response JSON ready")
		
			/// START
		fieldDataBytesJson,err := getJsonResponseError(errorGeneral, errorGeneralNbr)
		//////////    write the response (ERROR)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fieldDataBytesJson)	
		if(err!=nil){
			
		}
	
    } 
					
}

   
   
   //get


func handleDBGettokenizedcards(w http.ResponseWriter, r *http.Request) {
	defer func() {
		db.Connection.Close(nil)
	}()
    var errorGeneral string
    var errorGeneralNbr string
   	var requestData modelito.RequestTokenizedCards

    errorGeneral=""
    requestData, errorGeneral=obtainParmsGettokenizedcards(r,errorGeneral)
	////////////////////////////////////////////////validate parms
	/// START
    if errorGeneral=="" {

		errorGeneral,errorGeneralNbr= ProcessGettokenizedcards(w , requestData)
	}

    if errorGeneral!=""{
    	//send error response if any
    	//prepare an error JSON Response, if any
		log.Print("CZ   STEP Get the ERROR response JSON ready")
		
			/// START
		fieldDataBytesJson,err := getJsonResponseError(errorGeneral, errorGeneralNbr)
		//////////    write the response (ERROR)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fieldDataBytesJson)	
		if(err!=nil){
			
		}
	
    } 
					
}




// handleGeneratetokenized for receive and handle the request from client
func handleDBGeneratetokenized(w http.ResponseWriter, r *http.Request) {
	defer func() {
		db.Connection.Close(nil)
	}()
     var requestData modelito.RequestTokenized
     var errorGeneral string
     var errorGeneralNbr string
          
    errorGeneral=""
    requestData,errorGeneral =obtainParmsGeneratetokenized(r,errorGeneral)


	////////////////////////////////////////////////validate parms
	/// START
    
    if errorGeneral=="" {

		errorGeneral,errorGeneralNbr= ProcessGeneratetokenized(w , requestData)
	}

    if errorGeneral!=""{
    	//send error response if any
    	//prepare an error JSON Response, if any
		log.Print("CZ   STEP Get the ERROR response JSON ready")
		
			/// START
		fieldDataBytesJson,err := getJsonResponseError(errorGeneral, errorGeneralNbr)
		//////////    write the response (ERROR)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fieldDataBytesJson)	
		if(err!=nil){
			
		}
	
    } 
					
}


///////////////////////////////v4
///////////////////////////////v4



// v4handleDBProcesspayment  receive and handle the request from client, access DB
func v4handleDBPostProcesspayment(w http.ResponseWriter, r *http.Request) {
	defer func() {
		db.Connection.Close(nil)
	}()
    var errorGeneral string
    var errorGeneralNbr string
    var requestData modelito.RequestPayment
    
    errorGeneral=""
requestData,errorGeneral =obtainPostParmsProcessPayment(r,errorGeneral)  //logicrequest_post.go

	////////////////////////////////////////////////validate parms
	/// START
	////////////////////////////////////////////////validate parms
	/// START
    
    if errorGeneral=="" {

		errorGeneral,errorGeneralNbr= v4ProcessProcessPayment(w , requestData)    //logicbusiness.go 
	}

    if errorGeneral!=""{
    	//send error response if any
    	//prepare an error JSON Response, if any
		log.Print("CZ   STEP Get the ERROR response JSON ready")
		
			/// START
		fieldDataBytesJson,err := getJsonResponseError(errorGeneral, errorGeneralNbr)
		//////////    write the response (ERROR)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fieldDataBytesJson)	
		if(err!=nil){
			
		}
	
    } 
					
}

// v4handleDBProcesspayment  receive and handle the request from client, access DB
func v4handleDBProcesspayment(w http.ResponseWriter, r *http.Request) {
	defer func() {
		db.Connection.Close(nil)
	}()

    var errorGeneral string
    var	errorGeneralNbr string
    var requestData modelito.RequestPayment
    errorGeneral=""
requestData,errorGeneral =obtainParmsProcessPayment(r,errorGeneral)

	////////////////////////////////////////////////validate parms
	/// START
    
    if errorGeneral=="" {

		errorGeneral,errorGeneralNbr= v4ProcessProcessPayment(w , requestData)    //logicbusiness.go 
	}

    if errorGeneral!=""{
    	//send error response if any
    	//prepare an error JSON Response, if any
		log.Print("CZ   STEP Get the ERROR response JSON ready")
		
			/// START
		fieldDataBytesJson,err := getJsonResponseError(errorGeneral, errorGeneralNbr)
		//////////    write the response (ERROR)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fieldDataBytesJson)	
		if(err!=nil){
			
		}
	
    } 
					
}
