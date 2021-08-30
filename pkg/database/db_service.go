package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/reflectutil"

	"github.com/jmoiron/sqlx"
)

// DBService db service
type DBService struct {
	*sqlx.DB
}

// NewDBService instantiates new db service
func NewDBService(db *sqlx.DB) DBService {
	return DBService{db}
}

// MappingConfig provides addtional hints to DBService
type MappingConfig struct {
	TableName string
}

// Mappable structs returns sql specific meta data
type Mappable interface {
	Meta() MappingConfig
	PrimaryKey() int64
	SetPrimaryKey(int64)
}

// DBResult for capturing result information for inserts
type DBResult struct {
	id       int64
	affected int64
	err      error
}

// LastInsertId last inserted id
func (dbr DBResult) LastInsertId() (int64, error) {
	return dbr.id, dbr.err
}

// RowsAffected number of rows affected
func (dbr DBResult) RowsAffected() (int64, error) {
	return dbr.affected, dbr.err
}

// NamedExec wrapper for sqls.NamedExec.
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

// MustExec wrapper for sqlx.MustExec
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

// PrepareNamed wrapper for sqlx.PrepareNamed. However, it handles retrieval of new id
func (u DBService) PrepareNamed(sql string, arg interface{}) (sql.Result, *DBError) {
	db := u.DB
	result := DBResult{}

	tx := db.MustBegin()
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return result, &DBError{sql, err}
	}

	var id int
	err = stmt.Get(&id, arg)
	if err != nil {
		return result, &DBError{sql, err}
	}
	result.id = int64(id)

	err = tx.Commit()
	if err != nil {
		return result, &DBError{sql, err}
	}

	return result, nil
}

// Insert convenience function for insert statement
func (u DBService) Insert(tableName string, i interface{}) (sql.Result, *DBError) {
	dbTags := reflectutil.ListDBTagsFor(i)

	insertStatement := "INSERT INTO %s (%s) VALUES (%s) RETURNING id"

	columns := strings.Join(dbTags, ",")
	var dbTagsWithColon []string
	for _, tag := range dbTags {
		dbTagsWithColon = append(dbTagsWithColon, ":"+tag)
	}
	values := strings.Join(dbTagsWithColon, ",")

	insertStatement = fmt.Sprintf(insertStatement, tableName, columns, values)
	return u.PrepareNamed(insertStatement, i)
}

// Delete convenience function for delete statement
func (u DBService) Delete(sql string, id int64) *DBError {
	db := u.DB
	tx := db.MustBegin()

	tx.MustExec(sql, id)

	err := tx.Commit()
	if err != nil {
		return &DBError{sql, err}
	}

	return nil
}

// FindBy finds ide by given where clause
func (u UserService) FindBy(i interface{}, where string, args ...interface{}) *DBError {
	o, ok := i.(Mappable)
	if !ok {
		return &DBError{Query: "", Err: errors.New("object to save must implement Mappable")}
	}

	db := u.DB
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", o.Meta().TableName, where)
	err := db.Get(i, sql, args...)
	if err != nil {
		return &DBError{sql, err}
	}
	return nil
}

// FindByID finds by id
func (u DBService) FindByID(tableName string, id int64, i interface{}) (interface{}, *DBError) {
	db := u.DB
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tableName)
	err := db.Get(i, sql, id)
	if err != nil {
		return nil, &DBError{sql, err}
	}
	return i, nil
}

// FindByName finds by name
func (u DBService) FindByName(tableName string, name string, i interface{}) (interface{}, *DBError) {
	db := u.DB
	sql := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", tableName)
	err := db.Get(i, sql, name)
	if err != nil {
		return nil, &DBError{sql, err}
	}
	return i, nil
}

// Select executes sqlx.Select
func (u DBService) Select(dest interface{}, query string, args ...interface{}) (interface{}, *DBError) {
	db := u.DB
	err := db.Select(dest, query, args...)
	if err != nil {
		return nil, &DBError{query, err}
	}
	return dest, nil
}
