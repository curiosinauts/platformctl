package database

import "github.com/jmoiron/sqlx"

type DBService struct {
	*sqlx.DB
}

func NewDBService(db *sqlx.DB) DBService {
	return DBService{db}
}

func (u DBService) NamedExec(sql string, obj interface{}) *DBError {
	db := u.DB
	tx := db.MustBegin()

	_, err := tx.NamedExec(sql, obj)
	if err != nil {
		tx.Rollback()
		return &DBError{sql, err}
	}

	err = tx.Commit()
	if err != nil {
		return &DBError{sql, err}
	}

	return nil
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
