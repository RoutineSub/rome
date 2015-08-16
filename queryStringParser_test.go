package godapper

import (
	"testing"
	"database/sql"
)

func TestParseQueryStringReplacesPlaceHolder(t *testing.T) {
	t.Parallel()
	queryString := "select * from table where col = {a}"
	expectedString := "select * from table where col = ?"
	called := false
	mp := mockDatabase{func (s string) (*sql.Stmt, error){
			if s != expectedString {
				t.Fatalf("Exepcted %v but was %v",expectedString,s)
			}
			called = true
			return new(sql.Stmt),nil
		}}
	parseQueryString(&mp, queryString)
	if !called {
		t.Fatal("No call to database.Prepare")
	}
}

func TestParseQueryStringBraceLiteral(t *testing.T) {
	t.Parallel()
	queryString := "select * from {{}table} where col = 2"
	expectedString := "select * from {table} where col = 2"
	called := false
	mp := mockDatabase{func(s string) (*sql.Stmt, error){
			if s != expectedString {
				t.Fatalf("Exepcted %v but was %v",expectedString,s)
			}
			called = true
			return new(sql.Stmt), nil
		}}
	parseQueryString(&mp, queryString)
	if !called {
		t.Fatal("No call to database.Prepare")
	}
}

func TestParseQueryStringListOfKeys(t *testing.T) {
	t.Parallel()
	queryString := "select * from table where col = {hello world}"
	mp := mockDatabase{func (string) (*sql.Stmt, error) {
		return new(sql.Stmt), nil
	}}
	prepped, err := parseQueryString(&mp, queryString)
	if err != nil {
		t.Fatalf("Unexected error %v",err)
	}
	if len(prepped.argKeys) != 1 {
		t.Fatalf("Expected exactly 1 argument matcher but found %v", len(prepped.argKeys))
	}
	if prepped.argKeys[0] != "hello world" {
		t.Fatalf("Exepcted key value \"hello world\" but was %v", prepped.argKeys[0])
	}
}

func TestParseQueryStringMultipleKeys(t *testing.T) {
	t.Parallel()
	queryString := "select * from table where col = {1} and row = {2}"
	mp := mockDatabase{func (string) (*sql.Stmt, error){
			return new(sql.Stmt), nil
		}}
	prepped, err := parseQueryString(&mp, queryString)
	if err != nil {
		t.Fatalf("Unexected error %v", err)
	}
	if len(prepped.argKeys) != 2 {
		t.Fatalf("Expected exactly 2 argument matchers but found %v", len(prepped.argKeys))
	}
	if prepped.argKeys[0] != "1" {
		t.Fatalf("Exepected key 0 value \"1\" but was %v", prepped.argKeys[0])
	}
	if prepped.argKeys[1] != "2" {
		t.Fatalf("Expected key 1 value \"2\" but was %v", prepped.argKeys[1])
	}
}
