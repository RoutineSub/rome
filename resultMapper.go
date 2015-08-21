package dapper

import (
	"database/sql"
	"reflect"
	"errors"
	"fmt"
)

var (
	EOR = errors.New("End of Rows")
)

type Rows interface {
	Err() error
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
}

type valueScanner struct {
	field reflect.Value
}

func (v *valueScanner) Scan(src interface{}) error {
	if !v.field.CanSet() {
		return fmt.Errorf("Unable to set field %v", v.field)
	}
	v.field.Set(reflect.ValueOf(src))
	return nil
}

type mapScanner struct {
	key string
	m map[string]interface{}
}

func (m *mapScanner) Scan(src interface{}) error {
	m.m[m.key] = src
	return nil
}

type arrayScanner struct {
	index int
	array []interface{}
}

func (a *arrayScanner) Scan(src interface{}) error {
	if len(a.array) < a.index {
		return fmt.Errorf("Index %v is out of slice range", a.index)
	}
	a.array[a.index] = src
	return nil
}

type Result struct {
	Rows
}

// Lets you read the results of a SQL query as though they were and encoded stream
// Decode behaves differently depending on the type of d
//    If d is a map[string]interface{} the values are populated with keys set to the column names
//    If d is an []interface{} the position of the column number in the Result is used
//    If d is a struct or a pointer to a struct then the field names are used
func (r *Result) Decode(d interface{}) error {
	if r == nil {
		return errors.New("Cannot decode on nil result")
	}
	hasNext := r.Next()
	if !hasNext {
		return EOR
	}
	cols,_ := r.Columns()
	scanners := make([]interface{},len(cols))
	var findScanner func(interface{},int,string) sql.Scanner
	switch d.(type) {
	case map[string]interface{}:
		findScanner = func(d interface{}, i int, columnName string) sql.Scanner {
			return &mapScanner{columnName, d.(map[string]interface{})}
		}
	case []interface{}:
		findScanner = func (d interface{}, i int, columnName string) sql.Scanner {
			return &arrayScanner{i,d.([]interface{})}
		}
	default:
		findScanner = func(d interface{}, i int, columnName string) sql.Scanner {
			v := reflect.ValueOf(d)
			field := reflect.Indirect(v).FieldByName(columnName)
			return &valueScanner{field}
		}
	}
	for i,col := range(cols) {
		scanners[i] = findScanner(d, i, col)
	}
	err := r.Scan(scanners...)
	return err
}

func MapResult(rows Rows) Result {
	return Result{rows}
}



