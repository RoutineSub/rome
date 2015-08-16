package godapper

import (
	"database/sql"
)

func transact(s *sql.Stmt, tx transaction) (stmt,func(), error) {
	if tx != nil {
		s := tx.Stmt(s)
		return s, func() {
			s.Close()
		}, nil
	}
	return s, func(){}, nil
}

type parsedQueryString struct {
	statement *sql.Stmt
	argKeys []string
	t func(*sql.Stmt, transaction) (stmt, func(), error)
}

func parsed(statement *sql.Stmt, argKeys []string) *parsedQueryString {
	p := new(parsedQueryString)
	p.statement = statement
	p.argKeys = argKeys
	p.t = transact
	return p
}

func argV(argKeys []string, args map[string]interface{}) []interface{} {
	argv := make([]interface{},len(argKeys))
	for i,m := range(argKeys) {
		argv[i] = args[m]
	}
	return argv
}

func (parsed *parsedQueryString) execute(tx transaction, args map[string]interface{}) (int64, error) {
	s, d, err := parsed.t(parsed.statement, tx)
	if err != nil {
		return 0, err
	}
	defer d()
	result, err := s.Exec(argV(parsed.argKeys, args)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (parsed *parsedQueryString) query(tx transaction, args map[string]interface{}) (rows, error) {
	s, d, err := parsed.t(parsed.statement, tx)
	if err != nil {
		return nil, err
	}
	defer d()
	return s.Query(argV(parsed.argKeys, args)...)
}