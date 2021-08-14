package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	*DBService
}

func NewUserService(db *sqlx.DB) UserService {
	dbService := NewDBService(db)
	return UserService{&dbService}
}

func (u UserService) Add(user User) (sql.Result, *DBError) {
	sql := `
		INSERT INTO curiosity.user
			  (google_id, username, password, email, hashed_email, is_active, private_key, public_key) 
		VALUES 
			  (:google_id, :username, :password, :email, :hashed_email, :is_active, :private_key, :public_key)
		RETURNING
			  id
	`
	return u.PrepareNamed(sql, &user)
}

func (u UserService) AddUserRepo(userRepo UserRepo) (sql.Result, *DBError) {
	sql := `
		INSERT INTO user_repo
			  (uri, user_id) 
		VALUES 
			  (:uri, :user_id)
		RETURNING
			  id
	`
	return u.PrepareNamed(sql, &userRepo)
}

func (u UserService) AddIDERepo(ideRepo IDERepo) (sql.Result, *DBError) {
	sql := `
		INSERT INTO ide_repo
			  (user_ide_id, uri) 
		VALUES 
			  (:user_ide_id, :uri)
		RETURNING
			  id
	`
	return u.PrepareNamed(sql, &ideRepo)
}

func (u UserService) AddUserIDE(userIDE UserIDE) (sql.Result, *DBError) {
	sql := `
		INSERT INTO user_ide
			  (user_id, ide_id) 
		VALUES 
			  (:user_id, :ide_id)
		RETURNING
			  id
	`
	return u.PrepareNamed(sql, &userIDE)
}

func (u UserService) AddIDERuntimeInstall(ideRuntimeInstall IDERuntimeInstall) (sql.Result, *DBError) {
	sql := `
		INSERT INTO ide_runtime_install
			  (user_ide_id, runtime_install_id) 
		VALUES 
			  (:user_ide_id, :runtime_install_id)
		RETURNING
			  id
	`
	return u.PrepareNamed(sql, &ideRuntimeInstall)
}

func (u UserService) DeleteIDERuntimeInstall(id int64) *DBError {
	sql := `DELETE FROM ide_runtime_install WHERE id = $1`
	return u.DBService.Delete(sql, id)
}

func (u UserService) DeleteALLUserIDEsForUser(userID int64) *DBError {
	sql := `DELETE FROM user_ide WHERE id IN (SELECT id FROM user_ide WHERE user_id = $1)`
	return u.DBService.Delete(sql, userID)
}

func (u UserService) DeleteALLUserReposForUser(userID int64) *DBError {
	sql := `DELETE FROM user_repo WHERE id IN (SELECT id FROM user_repo WHERE user_id = $1)`
	return u.DBService.Delete(sql, userID)
}

func (u UserService) DeleteALLIDERuntimeInstallsForUser(userID int64) *DBError {
	sql := `DELETE FROM ide_runtime_install WHERE id IN (SELECT id FROM ide_runtime_install WHERE user_ide_id in (SELECT id FROM user_ide WHERE user_id = $1))`
	return u.DBService.Delete(sql, userID)
}

func (u UserService) DeleteALLIDEReposForUser(userID int64) *DBError {
	sql := `DELETE FROM ide_repo WHERE id IN (SELECT id FROM ide_repo WHERE user_ide_id in (SELECT id FROM user_ide WHERE user_id = $1))`
	return u.DBService.Delete(sql, userID)
}

func (u UserService) Get(id string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM curiosity.user WHERE user_id=$1"
	err := db.Get(&user, sql, id)
	if err != nil {
		return user, &DBError{sql, err}
	}
	return user, nil
}

func (u UserService) FindIDEByName(name string) (IDE, *DBError) {
	db := u.DB
	ide := IDE{}
	sql := "SELECT * FROM ide WHERE name=$1"
	err := db.Get(&ide, sql, name)
	if err != nil {
		return ide, &DBError{sql, err}
	}
	return ide, nil
}

func (u UserService) FindRuntimeInstallByName(name string) (RuntimeInstall, *DBError) {
	db := u.DB
	runtimeInstall := RuntimeInstall{}
	sql := "SELECT * FROM runtime_install WHERE name=$1"
	err := db.Get(&runtimeInstall, sql, name)
	if err != nil {
		return runtimeInstall, &DBError{sql, err}
	}
	return runtimeInstall, nil
}

func (u UserService) FindUserIDEsByUserID(userID int64) ([]int64, *DBError) {
	db := u.DB
	userIDEIDs := make([]int64, 0, 10)
	sql := "SELECT id FROM user_ide WHERE user_id=$1"
	err := db.Select(&userIDEIDs, sql, userID)
	if err != nil {
		return userIDEIDs, &DBError{sql, err}
	}
	return userIDEIDs, nil
}

func (u UserService) FindUserReposUserID(userID int64) ([]int64, *DBError) {
	db := u.DB
	userIDEIDs := make([]int64, 0, 10)
	sql := "SELECT id FROM user_ide WHERE user_id=$1"
	err := db.Select(&userIDEIDs, sql, userID)
	if err != nil {
		return userIDEIDs, &DBError{sql, err}
	}
	return userIDEIDs, nil
}

func (u UserService) FindIDEReposByUserID(userID int64) ([]int64, *DBError) {
	db := u.DB
	userIDEIDs := make([]int64, 0, 10)
	sql := "SELECT id FROM user_ide WHERE user_id=$1"
	err := db.Select(&userIDEIDs, sql, userID)
	if err != nil {
		return userIDEIDs, &DBError{sql, err}
	}
	return userIDEIDs, nil
}

func (u UserService) Delete(id int64) *DBError {
	db := u.DB
	sql := "DELETE FROM curiosity.user WHERE id=$1"
	_, err := db.Exec(sql, id)
	if err != nil {
		return &DBError{sql, err}
	}
	return nil
}

func (u UserService) FindByEmail(email string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM curiosity.user WHERE email=$1"
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

func (u UserService) UpdateProfile(user User) (sql.Result, *DBError) {
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