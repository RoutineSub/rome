package godapper

import "database/sql"

type mockTransaction struct {
	s func(*sql.Stmt) *sql.Stmt
}

func (m *mockTransaction) Commit() error {
	return nil
}

func (m *mockTransaction) Rollback() error {
	return nil
}

func (m *mockTransaction) Stmt(s *sql.Stmt) *sql.Stmt {
	return m.s(s)
}

type mockDatabase struct {
	prepare func(string) (*sql.Stmt, error)
}

func (m *mockDatabase) Begin() (*sql.Tx, error) {
	return nil, nil
}

func (m *mockDatabase) Prepare(queryString string) (*sql.Stmt, error) {
	return m.prepare(queryString)
}

type mockStatment struct {
	e func(args ...interface{}) (sql.Result, error)
	q func(args ...interface{}) (*sql.Rows, error)
}

func (m *mockStatment) Exec(args ...interface{}) (sql.Result, error) {
	return m.e(args...)
}

func (m *mockStatment) Query(args ...interface{}) (*sql.Rows, error) {
	return m.q(args...)
}

type mockResult struct {
	rows int64
}

func (m *mockResult) RowsAffected() (int64, error) {
	return m.rows, nil
}

func (m *mockResult) LastInsertId() (int64, error) {
	return 0, nil
}

type mockRows struct {
	e func()error
	pointer int
	columns []string
	values [][]interface{}
}

func (m *mockRows) Err() error {
	return m.e()
}

func (m *mockRows) Next() bool {
	m.pointer++
	return m.pointer < len(m.values)
}

func (m *mockRows) Columns() ([]string, error) {
	return m.columns, nil
}

func (m *mockRows) Scan(dest ...interface{}) error {
	for i := 0 ; i < len(m.values[m.pointer]) && i < len(dest); i++ {
		scan, ok := dest[i].(sql.Scanner)
		if ok {
			scan.Scan(m.values[m.pointer][i])
		} else {
			dest[i] = m.values[m.pointer][i]
		}
	}
	return nil
}

