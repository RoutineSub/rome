package godapper

import (
	"database/sql"
	"reflect"
)

type DB struct {
	db database
}

type Tx struct {
	tx transaction
	db database
}

var (
	globalCache = cache{make(map[string]*parsedQueryString), make(map[mapperKey]func(reflect.Value)sql.Scanner),
		parseQueryString, getScanner}
)

func Wrap(db *sql.DB) *DB {
	return &DB{db}
}

func (db *DB) Query(queryString string, args map[string]interface{}, dest interface{}) ([]interface{}, error) {
	parsed, err := globalCache.getCachedQuery(db.db, queryString)
	if err != nil {
		return nil, err
	}
	rows, err := parsed.query(nil, args)
	if err != nil {
		return nil, err
	}
	return mapResult(rows, reflect.TypeOf(dest), globalCache.getCachedMapper)
}

func (db *DB) Execute(queryString string, args map[string]interface{}) (int64,error) {
	parsed, err := globalCache.getCachedQuery(db.db, queryString)
	if err != nil {
		return 0, err
	}
	return parsed.execute(nil,args)
}

func (db *DB) Begin() (*Tx, error) {
	tx,err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx,db.db},nil
}

func (tx *Tx) Query(queryString string, args map[string]interface{}, dest interface{}) ([]interface{}, error) {
	parsed, err := globalCache.getCachedQuery(tx.db, queryString)
	if err != nil {
		return nil, err
	}
	rows, err := parsed.query(tx.tx, args)
	if err != nil {
		return nil, err
	}
	return mapResult(rows, reflect.TypeOf(dest), globalCache.getCachedMapper)
}

func (tx *Tx) Execute(queryString string, args map[string]interface{}) (int64, error) {
	parsed, err := globalCache.getCachedQuery(tx.db, queryString)
	if err != nil {
		return 0, err
	}
	return parsed.execute(tx.tx, args)
}

func (tx *Tx) Commit() error{
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error{
	return tx.tx.Rollback()
}