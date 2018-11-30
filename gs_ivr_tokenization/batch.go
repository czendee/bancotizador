package main

import (
	"banwire/services/gs_ivr_tokenization/db"
)

// BatchTest is a function only for test
func BatchTest() {

	defer func() {
		db.Connection.Close(nil)
	}()

	/*
		bla bla bla ... more code
	*/
}
