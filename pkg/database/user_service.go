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

// DeleteALLUserIDEsForUser deletes all user ides for given user
func (u UserService) DeleteALLUserIDEsForUser(userID int64) *DBError {
	query := `DELETE FROM user_ide WHERE id IN (SELECT id FROM user_ide WHERE user_id = $1)`
	return u.Delete(query, userID)
}

// DeleteALLUserReposForUser deletes all user repos for given user
func (u UserService) DeleteALLUserReposForUser(userID int64) *DBError {
	query := `DELETE FROM user_repo WHERE id IN (SELECT id FROM user_repo WHERE user_id = $1)`
	return u.DBService.Delete(query, userID)
}

// DeleteALLIDERuntimeInstallsForUser deletes all ide runtime installs for given user
func (u UserService) DeleteALLIDERuntimeInstallsForUser(userID int64) *DBError {
	query := `
		DELETE FROM ide_runtime_install WHERE id IN (
			SELECT id FROM ide_runtime_install WHERE user_ide_id in (
				SELECT id FROM user_ide WHERE user_id = $1
			)
		)`
	return u.DBService.Delete(query, userID)
}

// DeleteALLIDEReposForUser deletes all ide repos for given user
func (u UserService) DeleteALLIDEReposForUser(userID int64) *DBError {
	query := `
		DELETE FROM ide_repo WHERE id IN (
			SELECT id FROM ide_repo WHERE user_ide_id in (
				SELECT id FROM user_ide WHERE user_id = $1
			)
		)`
	return u.DBService.Delete(query, userID)
}

// FindUserIDEsByUserID find user ides by user id
func (u UserService) FindUserIDEsByUserID(userID int64) (*[]string, *DBError) {
	query := `
		SELECT name FROM ide WHERE id in (
			SELECT ide_id FROM user_ide WHERE user_id=$1
		)`
	i, dberr := u.Select(&[]string{}, query, userID)
	return i.(*[]string), dberr
}

// FindUserIDERepoURIsByUserAndIDE finds user ide repos by user and ide
func (u UserService) FindUserIDERepoURIsByUserAndIDE(username string, ide string) (*[]string, *DBError) {
	query := `
		SELECT uri FROM ide_repo WHERE user_ide_id = (
			SELECT 
					id as user_ide_id 
			FROM 
					user_ide 
			WHERE 
					user_id = (SELECT id as user_id FROM curiosity.user WHERE username = $1) AND
					ide_id  = (SELECT id ide_id     FROM ide            WHERE name     = $2)       
		)`
	i, dberr := u.Select(&[]string{}, query, username, ide)
	return i.(*[]string), dberr
}

// FindUserIDERuntimeInstallsByUserAndIDE finds runtime installs by user and ide
func (u UserService) FindUserIDERuntimeInstallsByUserAndIDE(username string, ide string) (*[]string, *DBError) {
	query := `
		SELECT script_body FROM runtime_install WHERE id in (
			SELECT runtime_install_id FROM ide_runtime_install WHERE user_ide_id = (
					SELECT 
							id as user_ide_id 
					FROM 
							user_ide 
					WHERE 
							user_id = (SELECT id as user_id FROM curiosity.user WHERE username = $1) AND
							ide_id  = (SELECT id ide_id     FROM ide            WHERE name     = $2)       
			)
		)`
	i, dberr := u.Select(&[]string{}, query, username, ide)
	return i.(*[]string), dberr
}

// FindUserIDERuntimeInstallNamesByUserAndIDE finds user install nemas by user and ide
func (u UserService) FindUserIDERuntimeInstallNamesByUserAndIDE(username string, ide string) (*[]string, *DBError) {
	query := `
		SELECT name FROM runtime_install WHERE id in (
			SELECT runtime_install_id FROM ide_runtime_install WHERE user_ide_id = (
					SELECT 
							id as user_ide_id 
					FROM 
							user_ide 
					WHERE 
							user_id = (SELECT id as user_id FROM curiosity.user WHERE username = $1) AND
							ide_id  = (SELECT id ide_id     FROM ide            WHERE name     = $2)       
			)
		) ORDER BY name ASC`
	i, dberr := u.Select(&[]string{}, query, username, ide)
	return i.(*[]string), dberr
}

// FindUserByGoogleID finds user by google id
func (u UserService) FindUserByGoogleID(googleIDHashed string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM curiosity.user WHERE google_id=$1"
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
			curiosity.user 
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
			curiosity.user 
		SET 
			google_id = :google_id
		WHERE 
			id = :id 
	`
	return u.NamedExec(sql, &user)
}
