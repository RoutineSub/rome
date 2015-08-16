package godapper

//Non exported interfaces for the database/sql package
//to enable testing and retargeting

import "database/sql"

type rows interface {
	Err() error
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
}

type stmt interface {
	Exec(args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
}

type transaction interface {
	Commit() error
	Rollback() error
	Stmt(*sql.Stmt) *sql.Stmt
}

type database interface {
	Begin() (*sql.Tx, error)
	Prepare(queryString string) (*sql.Stmt, error)
}

