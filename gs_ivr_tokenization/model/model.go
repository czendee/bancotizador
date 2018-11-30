package model

import (
	"encoding/json"
	"fmt"
	"time"

	"banwire/services/gs_ivr_tokenization/db"
)

// Time represents the time object
type Time struct {
	db.Time
}

// MarshalJSON returns the string json of time object like timestamp unix
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", t.UTC().Unix())), nil
}

// UnmarshalJSON returns the string json of time object like timestamp unix
func (t *Time) UnmarshalJSON(data []byte) error {
	var tt int64
	err := json.Unmarshal(data, &tt)

	if err == nil && tt > 0 {
		t.Time = db.Time{time.Unix(tt, 0), true}
	}

	return err
}

type Json interface{}
