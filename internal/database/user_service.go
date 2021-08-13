package database

import "github.com/jmoiron/sqlx"

type UserService struct {
	*DBService
}

func NewUserService(db *sqlx.DB) UserService {
	dbService := NewDBService(db)
	return UserService{&dbService}
}

func (u UserService) Add(user User) *DBError {
	sql := `
		INSERT INTO curiosity.user
			  (id,  google_id, username, password, email, hashed_email, is_active, private_key, public_key) 
		VALUES 
			  (:id, :google_id, :username, :password, :email, :hashed_email, :is_active, :private_key, :public_key)
	`
	return u.NamedExec(sql, &user)
}

func (u UserService) Get(id string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM users WHERE user_id=$1"
	err := db.Get(&user, sql, id)
	if err != nil {
		return user, &DBError{sql, err}
	}
	return user, nil
}

func (u UserService) Delete(id string) *DBError {
	db := u.DB
	sql := "DELETE FROM users WHERE user_id=$1"
	_, err := db.Exec(sql, id)
	if err != nil {
		return &DBError{sql, err}
	}
	return nil
}

func (u UserService) FindByEmail(email string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM users WHERE email=$1"
	err := db.Get(&user, sql, email)
	if err != nil {
		return user, &DBError{sql, err}
	}
	return user, nil
}

func (u UserService) List() ([]User, *DBError) {
	db := u.DB
	users := []User{}
	sql := "SELECT * FROM users"
	err := db.Select(&users, sql)
	if err != nil {
		return users, &DBError{sql, err}
	}
	return users, nil
}

func (u UserService) UpdateProfile(user User) *DBError {
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
