package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// UserService manages user data
type UserService struct {
	*DBService
}

// NewUserService instantiates new user service
func NewUserService(db *sqlx.DB) UserService {
	dbService := NewDBService(db)
	return UserService{&dbService}
}

// NewUserServiceWithOptions returns new user service with options
func NewUserServiceWithOptions(db *sqlx.DB, options ...DBOption) UserService {
	_dbs := &DBService{db, false}

	for _, option := range options {
		option(_dbs)
	}

	return UserService{_dbs}
}

// FindUserIDERuntimeInstallsByUsernameAndIDE finds user installs by user and ide
func (u UserService) FindAllRuntimeInstallsForUser(dest interface{}, username string) *DBError {
	query := `SELECT * FROM runtime_install WHERE name in ('tmux') ORDER BY name asc`
	_, dberr := u.Select(dest, query)
	return dberr
}

// FindUserByGoogleID finds user by google id
func (u UserService) FindUserByGoogleID(googleIDHashed string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM users WHERE google_id=$1"
	err := db.Get(&user, sql, googleIDHashed)
	if err != nil {
		return user, &DBError{sql, err}
	}
	return user, nil
}

// UpdateProfile updates user profile
func (u UserService) UpdateProfile(user User) (sql.Result, *DBError) {
	sql := `
		UPDATE 
			users
		SET 
			public_key_id = :public_key_id
		WHERE 
			id = :id 
	`
	return u.NamedExec(sql, &user)
}

// UpdateGoogleID updates user google id
func (u UserService) UpdateGoogleID(user User) (sql.Result, *DBError) {
	sql := `
		UPDATE 
			users 
		SET 
			google_id = :google_id
		WHERE 
			id = :id 
	`
	return u.NamedExec(sql, &user)
}
