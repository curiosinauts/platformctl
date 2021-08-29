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

// AddUser adds new user
func (u UserService) AddUser(user User) (sql.Result, *DBError) {
	return u.Insert("curiosity.user", &user)
}

// AddUserRepo adds new user repo
func (u UserService) AddUserRepo(userRepo UserRepo) (sql.Result, *DBError) {
	return u.Insert("user_repo", &userRepo)
}

// AddUserIDE adds new user ide
func (u UserService) AddUserIDE(userIDE UserIDE) (sql.Result, *DBError) {
	return u.Insert("user_ide", &userIDE)
}

// AddIDERepo adds new ide repo
func (u UserService) AddIDERepo(ideRepo IDERepo) (sql.Result, *DBError) {
	return u.Insert("ide_repo", &ideRepo)
}

// AddIDERuntimeInstall adds new ide runtime install
func (u UserService) AddIDERuntimeInstall(ideRuntimeInstall IDERuntimeInstall) (sql.Result, *DBError) {
	return u.Insert("ide_runtime_install", &ideRuntimeInstall)
}

// DeleteALLUserIDEsForUser deletes all user ides for given user
func (u UserService) DeleteALLUserIDEsForUser(userID int64) *DBError {
	query := `DELETE FROM user_ide WHERE id IN (SELECT id FROM user_ide WHERE user_id = $1)`
	return u.DBService.Delete(query, userID)
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

// FindByIDIDE finds ide by id
func (u UserService) FindByIDIDE(id int64) (*IDE, *DBError) {
	i, dberr := u.FindByID("ide", id, new(IDE))
	return i.(*IDE), dberr
}

// FindByNameIDE finds ide by name
func (u UserService) FindByNameIDE(name string) (*IDE, *DBError) {
	i, dberr := u.FindByName("ide", name, new(IDE))
	return i.(*IDE), dberr
}

// FindByIDRuntimeInstall finds runtime install by id
func (u UserService) FindByIDRuntimeInstall(id int64) (*RuntimeInstall, *DBError) {
	i, dberr := u.FindByID("runtime_install", id, new(RuntimeInstall))
	return i.(*RuntimeInstall), dberr
}

// FindByNameRuntimeInstall finds runtime install by name
func (u UserService) FindByNameRuntimeInstall(name string) (*RuntimeInstall, *DBError) {
	i, dberr := u.FindByName("runtime_install", name, new(RuntimeInstall))
	return i.(*RuntimeInstall), dberr
}

// FindUserIDEsByUser finds UserIDE by user
func (u UserService) FindUserIDEsByUser(user User) (*[]UserIDE, *DBError) {
	i, dberr := u.Select(&[]UserIDE{}, "SELECT * FROM user_ide WHERE user_id=$1", user.ID)
	return i.(*[]UserIDE), dberr
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
		)`
	i, dberr := u.Select(&[]string{}, query, username, ide)
	return i.(*[]string), dberr
}

func (u UserService) FindIDERuntimeInstallsByUserIDE(userIDE UserIDE) ([]IDERuntimeInstall, *DBError) {
	db := u.DB
	ideRuntimeInstalls := []IDERuntimeInstall{}
	sql := `SELECT * FROM ide_runtime_install WHERE user_ide_id = $1`
	err := db.Select(&ideRuntimeInstalls, sql, userIDE.ID)
	if err != nil {
		return ideRuntimeInstalls, &DBError{sql, err}
	}
	return ideRuntimeInstalls, nil
}

func (u UserService) DeleteUser(id int64) *DBError {
	db := u.DB
	sql := "DELETE FROM curiosity.user WHERE id=$1"
	_, err := db.Exec(sql, id)
	if err != nil {
		return &DBError{sql, err}
	}
	return nil
}

func (u UserService) FindUserByHashedEmail(hashedEmail string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM curiosity.user WHERE hashed_email=$1"
	err := db.Get(&user, sql, hashedEmail)
	if err != nil {
		return user, &DBError{sql, err}
	}
	return user, nil
}

func (u UserService) FindUserByUsername(username string) (User, *DBError) {
	db := u.DB
	user := User{}
	sql := "SELECT * FROM curiosity.user WHERE username=$1"
	err := db.Get(&user, sql, username)
	if err != nil {
		return user, &DBError{sql, err}
	}
	return user, nil
}

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

func (u UserService) List() ([]User, *DBError) {
	db := u.DB
	users := []User{}
	sql := "SELECT * FROM curiosity.user"
	err := db.Select(&users, sql)
	if err != nil {
		return users, &DBError{sql, err}
	}
	return users, nil
}

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
