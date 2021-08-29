package database

import (
	"encoding/json"
	"log"
)

// DBError database error
type DBError struct {
	Query string `json:"query"`
	Err   error  `json:"error"`
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// The error interface implementation, which formats to a JSON object string.
func (e *DBError) Error() string {
	marshaled, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(marshaled)
}

// Log logs the database transaction with query and error
func (e *DBError) Log(tx string) {
	log.Printf("%s sql.excute.errored: %#v, sql: %s", tx, e.Err, e.Query)
}
