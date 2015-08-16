package godapper

import (
	"testing"
	"database/sql"
)

func TestTransactWrapsInTransactionWhenPresent(t *testing.T){
	t.Parallel()
	called := false
	expected := new(sql.Stmt)
	mt := mockTransaction{func (s *sql.Stmt) *sql.Stmt{
			if s != expected {
				t.Fatal("Expected %v but was %v", expected, s)
			}
			called = true
			return s
		}}
	transact(expected, &mt)
	if !called {
		t.Fatal("No call to transaction.Stmt")
	}
}

func TestArgVMapsValues(t *testing.T){
	t.Parallel()
	values:= map[string]interface{}{"test":10,"other":"cat"}
	args:= []string{"test","other"}
	argv := argV(args, values)
	if len(argv) != 2 {
		t.Fatalf("Expected length 2 but was %v",len(argv))
	}
	if argv[0] != 10 {
		t.Fatalf("Expected argument 0 to be 10 but was %v", argv[0])
	}
	if argv[1] != "cat" {
		t.Fatalf("Expected argument 1 to be \"cat\" but was %v", argv[1])
	}
}

func TestExecuteCallsStatementExecuteWithArgs(t *testing.T){
	t.Parallel()
	called := false
	ms := mockStatment{func(args ...interface{}) (sql.Result, error){
			if len(args) != 1 {
				t.Fatalf("Exected 1 argument but was %v", len(args))
			}
			if args[0] != 10 {
				t.Fatalf("Expected argument value 10 but was %v", args[0])
			}
			called = true
			return &mockResult{0}, nil
		}, func(...interface{}) (*sql.Rows, error){
			return nil, nil
		}}
	p := parsedQueryString{new(sql.Stmt),[]string{"test"}, func(*sql.Stmt, transaction) (stmt, func(), error){
			return &ms, func(){}, nil
	}}
	p.execute(nil, map[string]interface{}{"test":10})
	if !called {
		t.Fatal("No call to stmt.Exec")
	}
}

func TestExecuteReturnsResultRows(t *testing.T) {
	t.Parallel()
	ms := mockStatment{func(...interface{}) (sql.Result, error){
			return &mockResult{102}, nil
		}, func(...interface{}) (*sql.Rows, error){
			return nil, nil
		}}
	p := parsedQueryString{new(sql.Stmt),[]string{}, func(*sql.Stmt, transaction) (stmt, func(), error){
			return &ms, func(){}, nil
		}}
	rows, _ := p.execute(nil, map[string]interface{}{})
	if rows != 102 {
		t.Fatalf("Expected 101 but was %v", rows)
	}
}

func TestExecuteCallsDeferedFunction(t *testing.T) {
	t.Parallel()
	called := false
	ms := mockStatment{func(...interface{}) (sql.Result, error){
		return &mockResult{0}, nil
	}, func(...interface{}) (*sql.Rows, error){
		return nil,nil
	}}
	trans := func(*sql.Stmt, transaction) (stmt, func(), error){
		return &ms, func(){
				called = true
			}, nil
	}
	p := parsedQueryString{new(sql.Stmt),[]string{},trans}
	p.execute(nil,map[string]interface{}{})
	if !called {
		t.Fatal("Expected call of deffered function")
	}
}

func TestQueryCallsStatementQueryWithArgs(t *testing.T){
	t.Parallel()
	called := false
	ms := mockStatment{func(...interface{})(sql.Result, error){
		return nil,nil
	}, func(args ...interface{}) (*sql.Rows, error){
		if len(args) != 1 {
			t.Fatalf("Expected exactly 1 argument but found %v", len(args))
		}
		if args[0] != 10 {
			t.Fatalf("Expected argument value 10 but was %v", args[0])
		}
		called = true;
		return new(sql.Rows), nil
	}}
	p := parsedQueryString{new(sql.Stmt), []string{"test"}, func(*sql.Stmt, transaction) (stmt, func(), error){
			return &ms, func(){}, nil
		}}
	p.query(nil,map[string]interface{}{"test":10})
	if !called {
		t.Fatal("Expected at call to stmt.Query")
	}
}