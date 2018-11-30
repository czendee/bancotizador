package net

import (
	"log"
	"net/http"
	"strconv"
	"sync"
//	"encoding/json"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
//	"time"

)

var lock = &sync.RWMutex{}

// HTTPMethod is the list of supported HTTP methods.
type HTTPMethod string

const (
	POST   HTTPMethod = "POST"
	GET    HTTPMethod = "GET"
	HEAD   HTTPMethod = "HEAD"
	DELETE HTTPMethod = "DELETE"
)

// Intializes the router instance
var router = mux.Router{}

// GetRouter returns the router object
func GetRouter() *mux.Router {
	return &router

}

var NoMatch func(http.ResponseWriter, *http.Request)
var HandleDefer func(http.ResponseWriter, *http.Request)

// CrossDomainRequest sets the headers in the response for allow crossdomain
func CrossDomainRequest(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

	if r.Method == "OPTIONS" {
		return false
	}

	return true
}

// Handler handles and redirects the requets using the router
func handler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var match mux.RouteMatch
//	var result string
	
    log.Print("CZ    handler Listening test 01")
/*    result ="NOK"
    if(r.URL.Path !="/"){
        log.Print("CZ    handler Listening test 02")
    	if(r.URL.Path =="/gettokenizedcards"){
    		    log.Print("CZ    handler Listening test obtienetarjetastokenizadas")
    		    
    		    	err := r.ParseForm()
					if err != nil {
						panic(err)
					}
					v := r.Form
					email := v.Get("email")
					idcomercio := v.Get("idcomercio")
	
		        result ="OK gettokenizedcards"+email+idcomercio
				//response
    		    log.Print("CZ    handler Listening test gettokenizedcards  2")					
				
//					fieldDataBytesJson,err := getJsonResponse()
//					fieldDataBytesJson,err := getJsonResponse01()
					
					fieldDataBytesJson,err := getJsonResponse02()
					
					log.Print("CZ    handler Listening test gettokenizedcards  3")	
	
		        result ="OK gettokenizedcards"+email+idcomercio+"resultado"
		        
		        w.Header().Set("Content-Type", "application/json")
		         w.Write(fieldDataBytesJson)

	    }else     	if(r.URL.Path =="/gettokenizedcardsjson"){
    		    log.Print("CZ    handler Listening test gettokenizedcardsjson")
    		    // process the request 
    		    
				decoder := json.NewDecoder(r.Body)
			
				var data myDataService01
				err := decoder.Decode(&data)
				if err != nil {
					panic(err)
				}
			
				email := data.Email
				idcomercio := data.Idcomercio
				
				//logic
				
				
				//response
    		    log.Print("CZ    handler Listening test gettokenizedcardsjson  2")					
					fieldDataBytesJson,err := getJsonResponse02()
    		    log.Print("CZ    handler Listening test gettokenizedcardsjson  3")	
	
		        result ="OK gettokenizedcardsjson:"+email+idcomercio+"resultado"
		        
		        w.Header().Set("Content-Type", "application/json")
		         w.Write(fieldDataBytesJson)
				if !router.Match(r, &match) {
					if NoMatch != nil {
						NoMatch(w, r)
					}
				} else {
					router.ServeHTTP(w, r)
				}
		        return

	    }else if(r.URL.Path =="/processpayment"){
    		    log.Print("CZ    handler Listening test realizarpago")
		        result ="OK processpayment"

	    }else if(r.URL.Path =="/generatetokenized"){
		        result ="OK generatetokenized"
	    }else{
	    	log.Print("CZ    handler Listening test not found")
		   http.NotFound(w,r)
		   return
	    }
            log.Print("CZ    handler Listening test 04")

	}
	    log.Print("CZ    handler Listening test 05")
	    
	arr := []byte("<html><body>"+ result+"</body></html>")  

    w.Write(arr)
    
    */
	if !router.Match(r, &match) {
		if NoMatch != nil {
			NoMatch(w, r)
		}
	} else {
		router.ServeHTTP(w, r)
	}
}

// HTTPListener initializes the HTTP listener
func HTTPListener(listening string) (err error) {
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(listening, nil)
	    log.Print("CZ    HTTPListener")
	if err != nil {
		log.Panic("Error initializing HTTP listener: %s", err.Error())
	}

	return
}

// ResponseWriter works like proxy of the Http ResponseWriter object
type ResponseWriter struct {
	Status int
	RW     http.ResponseWriter
	Body   []byte

	WriteCallback func([]byte) int
}

// Headers returns the header object
func (rw *ResponseWriter) Header() http.Header {
	return rw.RW.Header()
}

// Write redirects and set the message HTTP response message
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.Body = b

	if status := rw.WriteCallback(b); status > 0 {
		rw.RW.WriteHeader(status)
	} else if rw.Status > 0 {
		rw.RW.WriteHeader(rw.Status)
	}

	return rw.RW.Write(b)
}

// WriteHeader redirects and set the HTTP response code
func (rw *ResponseWriter) WriteHeader(i int) {
	if i != 0 {
		rw.Status = i
	}
}

// HttpInterface is a interface that support the HTTP methods
type HTTPInterface interface {
	HTTPStatus() int
}

// ResponseWriterJSON sets the heads for a JSON response message
// and returns the ResponseWriter object
func ResponseWriterJSON(w http.ResponseWriter, obj ...interface{}) http.ResponseWriter {
	var rw = &ResponseWriter{
		RW: w,
		WriteCallback: func(body []byte) int {
			length := strconv.Itoa(len(body))
			w.Header().Set("Content-Length", length)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if len(obj) > 0 {
				if i, ok := obj[0].(HTTPInterface); ok {
					return i.HTTPStatus()
				}
			}
			return 0
		},
	}

	return rw
}

// Handle returns negroni object that waps the handler function for the router object
func Handle(h http.HandlerFunc) *negroni.Negroni {
	return negroni.New(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		defer func() {
			if HandleDefer != nil {
				HandleDefer(w, r)
			}
		}()

		next(w, r)
	}),
		negroni.Wrap(http.HandlerFunc(h)),
	)
}
/*

type myDataService01 struct {
	Email string            `json:"Email_client"`
	Idcomercio  string      `json:"Idcomercio"`
}



type AutoGenerated struct {
	StatusMessage string       `json:"Status_message"`
	Tarjetas       []TarjetaData `json:"Tarjetas_data"`
}

type TarjetaData struct {
	Date time.Time `json:"Tarjeta_date"`
	Token  string   `json:"Tarjeta_token"`
    LastDigits  string   `json:"Tarjeta_last"`
	Marca  string   `json:"Tarjeta_marca"`
	Vigencia  string   `json:"Tarjeta_vigencia"`
}

func getJsonResponse02()([]byte, error) {
	
	mainStruct := AutoGenerated{StatusMessage: "Success"}

	for i := 0; i < 5; i++ {
		w := TarjetaData{time.Now(), "G#$$%ytoywteouwytr","1234","VISA","2501"}
		mainStruct.Tarjetas = append(mainStruct.Tarjetas, w)
	}


	return json.MarshalIndent(mainStruct, "", "  ")
}
*/