package database

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type UserService struct {
	*sqlx.DB
}

func CurrentTimeInSeconds() int64 {
	return time.Now().Unix()
}

type DBOpError struct {
	Query string `json:"query"`
	Err   error  `json:"error"`
}

func (e *DBOpError) Unwrap() error {
	return e.Err
}

// The error interface implementation, which formats to a JSON object string.
func (e *DBOpError) Error() string {
	marshaled, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(marshaled)
}

func (e *DBOpError) Log(tx string) {
	log.Printf("%s sql.excute.errored: %#v, sql: %s", tx, e.Err, e.Query)
}

func (u UserService) NamedExec(sql string, obj interface{}) *DBOpError {
	db := u.DB
	tx := db.MustBegin()

	_, err := tx.NamedExec(sql, obj)
	if err != nil {
		tx.Rollback()
		return &DBOpError{sql, err}
	}

	err = tx.Commit()
	if err != nil {
		return &DBOpError{sql, err}
	}

	return nil
}

func (u UserService) MustExec(sql string, args ...interface{}) *DBOpError {
	db := u.DB
	tx := db.MustBegin()

	tx.MustExec(sql, args...)

	err := tx.Commit()
	if err != nil {
		return &DBOpError{sql, err}
	}

	return nil
}

func (u UserService) Register(user User) *DBOpError {
	sql := `
		INSERT INTO users
			  (user_id,  platform_name,  email,   passhash,  firstname, lastname,  created_date) 
		VALUES 
			  (:user_id, :platform_name, :email, :passhash, :firstname, :lastname, :created_date)
	`
	return u.NamedExec(sql, &user)
}

func (u UserService) Get(id string) (User, *DBOpError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM users WHERE user_id=$1"
	err := db.Get(&user, sql, id)
	if err != nil {
		return user, &DBOpError{sql, err}
	}
	return user, nil
}

func (u UserService) Delete(id string) *DBOpError {
	db := u.DB
	sql := "DELETE FROM users WHERE user_id=$1"
	_, err := db.Exec(sql, id)
	if err != nil {
		return &DBOpError{sql, err}
	}
	return nil
}

func (u UserService) FindByEmail(email string) (User, *DBOpError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM users WHERE email=$1"
	err := db.Get(&user, sql, email)
	if err != nil {
		return user, &DBOpError{sql, err}
	}
	return user, nil
}

func (u UserService) List() ([]User, *DBOpError) {
	db := u.DB
	users := []User{}
	sql := "SELECT * FROM users"
	err := db.Select(&users, sql)
	if err != nil {
		return users, &DBOpError{sql, err}
	}
	return users, nil
}

func (u UserService) UpdateProfile(user User) *DBOpError {
	sql := `
		UPDATE 
			users 
		SET 
			firstname        = :firstname,  
		    lastname         = :lastname,
			email            = :email,
			passhash         = :passhash,
			totp_enabled     = :totp_enabled,
            webauthn_enabled = :webauthn_enabled
		WHERE 
			user_id = :user_id 
	`
	return u.NamedExec(sql, &user)
}
