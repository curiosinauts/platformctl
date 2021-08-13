package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DBService struct {
	*sqlx.DB
}

func NewDBService(db *sqlx.DB) DBService {
	return DBService{db}
}

type DBResult struct {
	id       int64
	affected int64
}

func (dbr DBResult) LastInsertId() (int64, error) {
	return dbr.id, nil
}

func (dbr DBResult) RowsAffected() (int64, error) {
	return dbr.affected, nil
}

func (u DBService) NamedExec(sql string, obj interface{}) (sql.Result, *DBError) {
	db := u.DB
	tx := db.MustBegin()

	result, err := tx.NamedExec(sql, obj)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{sql, err}
	}

	err = tx.Commit()
	if err != nil {
		return nil, &DBError{sql, err}
	}

	return result, nil
}

func (u DBService) MustExec(sql string, args ...interface{}) *DBError {
	db := u.DB
	tx := db.MustBegin()

	tx.MustExec(sql, args...)

	err := tx.Commit()
	if err != nil {
		return &DBError{sql, err}
	}

	return nil
}

func (u DBService) PrepareNamed(sql string, args ...interface{}) (sql.Result, *DBError) {
	db := u.DB
	result := DBResult{}

	sql, _, err := db.BindNamed(sql, args)
	if err != nil {
		panic(err)
	}

	//tx := db.MustBegin()
	stmt, err := db.PrepareNamed(sql)

	var id int
	err = stmt.Get(&id, args)
	result.id = int64(id)

	//err = tx.Commit()
	if err != nil {
		return result, &DBError{sql, err}
	}

	return result, nil
}
