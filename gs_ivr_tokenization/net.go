package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	//"banwire/services/gs_ivr_tokenization/db"
	//"banwire/services/gs_ivr_tokenization/model/pgsql"
	"banwire/services/gs_ivr_tokenization/model"
	"banwire/services/gs_ivr_tokenization/net"

	"github.com/codegangsta/negroni"
	//"github.com/gorilla/context"
)



//init initializes
func init() {
	net.NoMatch = func(w http.ResponseWriter, r *http.Request) {
		var err = model.NotFoundError("Method not found")
		rw := net.ResponseWriterJSON(w, err)
		response, _ := json.Marshal(err)
		rw.Write(response)
	}
}

// AccessInterface is an interface that supports all methods for the access objects
type AccessInterface interface {
	Access(*http.Request) (bool, error)
}

// AccessApiKey handles the access and validation by api key
type AccessApiKey struct {
	ApiKeyValidation bool

	Permission []string

	CrossDomain bool
}

// SetApiKeyValidation sets the ApiKeyValidation parameter value
func (a *AccessApiKey) SetApiKeyValidation(apiKeyValidation bool) *AccessApiKey {
	a.ApiKeyValidation = apiKeyValidation
	return a
}

// SetCrossDomain sets the CrossDomain parameter value
func (a *AccessApiKey) SetCrossDomain(crossDomain bool) *AccessApiKey {
	a.CrossDomain = crossDomain
	return a
}

// Access returns the result of the validation
func (a *AccessApiKey) Access(r *http.Request) (ok bool, err error) {
	/*
		var permission = a.Permission
		var apiKeyValidation = a.ApiKeyValidation

		var api_key *model.ApiKey

		_pg := pgsql.Db{
			db.Connection.GetPgDb(r, false),
		}

		var key, secret string

		if apiKeyValidation {
			key, secret, ok = r.BasicAuth()
		}

		if secret == "_crossdomain_" {
			secret = ""
		}
		if secret == "" && a.CrossDomain {
			secret = "_crossdomain_"
		}

		if ok {
			a, e := _pg.SecurityApiKeyAuthFunction(key, secret)
			api_key = &a

			if e != nil {
				if _e, _ok := e.(*pgsql.Error); _ok {
					switch _e.Type {
					case "auth_invalid":
						err = model.UnauthorizedError("Authentication failed")
						return
					}
					log.Print(_e.Type)
				}
				log.Panic(e)
			} else if permission != nil && !api_key.IsAllowed(permission, true) {
				err = model.UnauthorizedError("Access denied. The Api Key does not have sufficient permissions to access")
				return
			} else {
				ok = true
				context.Set(r, model.APIKEY, api_key)
			}
		} else {
			err = model.UnauthorizedError("Authentication is required")
			return
		}

	*/
	return
}

// NewAccessApiKey initializes and returns the access object by api key (AccessApiKey)
func NewAccessApiKey(apiKeyValidation bool, permission ...string) *AccessApiKey {
	return &AccessApiKey{
		ApiKeyValidation: apiKeyValidation,
		Permission:       permission,
	}
}

// netHandle returns negroni object that waps the handler function for the router object
func netHandle(h http.HandlerFunc, access AccessInterface) *negroni.Negroni {
	return negroni.New(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var err error

		defer func() {
			if err != nil {
				rw := net.ResponseWriterJSON(w, err)
				response, _ := json.Marshal(err)
				rw.Write(response)
			}

			if _r := recover(); _r != nil {
				log.Print(_r)
				ip := IdentifyPanic()
				log.Print(ip)

				err, _ := json.Marshal(model.InvalidRequestError("The Api has experienced an internal error."))
				rw := net.ResponseWriterJSON(w)
				rw.WriteHeader(500)
				rw.Write(err)

				if strings.Contains(ip, "jackc/pgx") {
					log.Fatal("exit")
				}
			}
		}()

		var ok bool = true

		if !net.CrossDomainRequest(w, r) {
			w.Write([]byte{})
			return
		}

		if access != nil {
			if ok, err = access.Access(r); err != nil {
				return
			}
		}

		if ok {
			next(w, r)
		}
	}),
		negroni.Wrap(http.HandlerFunc(h)),
	)
}

func IdentifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	paths := []string{}

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			switch {
			case name != "":
				paths = append(paths, fmt.Sprintf("%v:%v", name, line))
			case file != "":
				paths = append(paths, fmt.Sprintf("%v:%v", file, line))
			}
		}
	}

	return strings.Join(paths, "\n")
}
