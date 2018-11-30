package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"banwire/services/gs_ivr_tokenization/db"
	"banwire/services/gs_ivr_tokenization/path"
)

// Initializes the Config variables
var Config config

var HTTPListen string
var configFile string

var RunMode string



func init() {
	flag.StringVar(&RunMode, "mode", "api", "Service mode run (Options: api, batch)")
	flag.StringVar(&HTTPListen, "http", ":8090", "Path where HTTP Server will listening")
	flag.StringVar(&configFile, "config", "./conf/config.json", "Path of configuration file")
//	flag.StringVar(&configFile, "config", "config.json", "Path of configuration file")

	configFile = path.RelativePath(configFile)
}


// loadConfiguration loads the configuration file
func LoadConfiguration() {
	log.Print("Loading configuration...")

	if d, e := ioutil.ReadFile(configFile); e == nil {
		e := json.Unmarshal(d, &Config)
		if e != nil {
//			log.Panicf("Error in unmarshalling the configuration file: %s", e.Error())
            log.Print("Configuration was not was loaded!")
		}else{
			log.Print("Configuration was YESS was loaded!")
		}

        
		log.Print("Configuration was loaded!(check previo")
	} else {
		//log.Panicf("Error in opening the configuration file: %s", e)
			log.Print("Error in opening the configuration file %s", e)
	}
}

// config is the configuration structure object
type config struct {
	Database configDatabase `json:"database"`
}

// configDatabase is the database structure object for configuration
type configDatabase struct{}

// UnmarshalJSON handles desearialization of configDatabase
// and loads the database connections
func (c *configDatabase) UnmarshalJSON(data []byte) error {

	var cc = []struct {
		Drive string                   `json:"drive"`
		Nodes []map[string]interface{} `json:"nodes"`
	}{}
			log.Print("UnmarshalJSON 01!")
	err := json.Unmarshal(data, &cc)
	if err != nil {
		return err
	}
			log.Print("UnmarshalJSON 02!")
	for _, d := range cc {
					log.Print("UnmarshalJSON 03.!"+d.Drive)
		switch d.Drive {
		case "postgresql":
			log.Print("UnmarshalJSON 04!")
			for _, n := range d.Nodes {
							log.Print("UnmarshalJSON 05!")
				if active, _ := n["active"].(bool); active {
					host, _ := n["host"].(string)
					port, _ := n["port"].(float64)
					_db, _ := n["db"].(string)
					user, _ := n["user"].(string)
					pass, _ := n["password"].(string)

                    
					if e := db.Connection.Set(db.NewPgDb(host, int(port), _db, user, pass)); e == nil {
						log.Print("---- The postgresql database was loaded"+host)
						log.Print("---- The postgresql database was loaded"+_db)
					} else {
						return e
					}
				}
							log.Print("UnmarshalJSON 06!")
			}

			break
		}
			log.Print("UnmarshalJSON 07!")
	}

	return nil
}
