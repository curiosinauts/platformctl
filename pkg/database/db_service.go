package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/reflectutil"

	"github.com/jmoiron/sqlx"
)

// DBService db service
type DBService struct {
	*sqlx.DB
	Debug bool
}

// NewDBService instantiates new db service
func NewDBService(db *sqlx.DB) DBService {
	return DBService{db, false}
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
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", o.Meta().TableName, where)
	if u.Debug {
		PrintQuery(query, args...)
	}
	err := db.Get(i, query, args...)
	if err != nil {
		return &DBError{query, err}
	}
	if u.Debug {
		PrintResult(i)
	}
	return nil
}

func PrintQuery(query string, args ...interface{}) {
	fmt.Println()
	fmt.Println(query, args)
}

func PrintResult(i interface{}) {
	fmt.Println(i)
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
	if u.Debug {
		PrintQuery(query, args...)
	}
	err := db.Select(dest, query, args...)
	if err != nil {
		return nil, &DBError{query, err}
	}
	if u.Debug {
		PrintResult(dest)
	}
	return dest, nil
}

// Save adds new data
func (u DBService) Save(i interface{}) *DBError {
	o, ok := i.(Mappable)
	if !ok {
		return &DBError{Query: "", Err: errors.New("object to save must implement Mappable")}
	}
	result, dberr := u.Insert(o.Meta().TableName, i)
	if dberr != nil {
		return dberr
	}

	id, err := result.LastInsertId()
	if err != nil {
		return &DBError{Query: "", Err: errors.New("error retrieving last insert id")}
	}
	o.SetPrimaryKey(id)
	return nil
}

// Del deletes given data
func (u DBService) Del(i interface{}) *DBError {
	o, ok := i.(Mappable)
	if !ok {
		return &DBError{Query: "", Err: errors.New("Object to save must implement Mappable")}
	}
	query := fmt.Sprintf("DELETE from %s WHERE id = $1", o.Meta().TableName)

	return u.Delete(query, o.PrimaryKey())
}

// List all rows of given entity
func (u DBService) List(tableName string, dest interface{}) *DBError {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	db := u.DB
	if u.Debug {
		PrintQuery(query)
	}
	err := db.Select(dest, query)
	if err != nil {
		return &DBError{query, err}
	}
	if u.Debug {
		PrintResult(dest)
	}
	return nil
}

// ListBy lists all rows of given entity with where clause
func (u DBService) ListBy(tableName string, dest interface{}, where string, args ...interface{}) *DBError {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, where)

	db := u.DB
	if u.Debug {
		PrintQuery(query, args...)
	}
	err := db.Select(dest, query, args...)
	if err != nil {
		return &DBError{query, err}
	}
	if u.Debug {
		PrintResult(dest)
	}
	return nil
}

func GetMappingConfigFromSlicePointer(dest interface{}) *MappingConfig {
	items := reflect.ValueOf(dest)
	if items.Kind() == reflect.Ptr && items.Elem().Kind() == reflect.Slice {
		itemsDeref := reflect.Indirect(items)
		for i := 0; i < itemsDeref.Len(); i++ {
			//fmt.Println("index", i)
			item := itemsDeref.Index(i)
			if item.Kind() == reflect.Struct {
				s := reflect.Indirect(item)
				for j := 0; j < s.NumField(); j++ {
					fmt.Println(s.Type().Field(j).Name, ":", s.Field(j).Interface())
				}
				for k := 0; k < s.NumMethod(); k++ {
					m := s.Method(k)
					fmt.Println("method: ", m, reflect.Indirect(m))
				}
				//val := v.MethodByName("Meta").Call([]reflect.Value{})
				//fmt.Println(val)
				//if !ok {
				//	fmt.Println("Mappable type assertion failed")
				//	return nil
				//}
				//meta := mappable.Meta()
				//return nil
			}
		}
	}
	return nil
}
