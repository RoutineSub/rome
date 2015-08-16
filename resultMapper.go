package godapper

import (
	"database/sql"
	"reflect"
	"fmt"
	"errors"
)

type valueScanner struct {
	value reflect.Value
	index []int
}

func (v *valueScanner) Scan(src interface{}) error {
	field := v.value.FieldByIndex(v.index)
	if !field.CanSet() {
		return errors.New("")
	}
	field.Set(reflect.ValueOf(src))
	return nil
}

func getScanner(t reflect.Type, columnName string) (func(reflect.Value)sql.Scanner, error) {
	field,found := t.FieldByName(columnName)
	if !found {
		return nil, fmt.Errorf("No field named %v in type %v",columnName,t.Name())
	}
	index := field.Index
	return func(v reflect.Value) sql.Scanner {
		return &valueScanner{v, index}
	}, nil
}

func mapScanners (scanners []func(reflect.Value)sql.Scanner, rowValue reflect.Value) []interface{} {
	rSlice := make([]interface{},len(scanners))
	for i,scanner := range(scanners) {
		rSlice[i] = scanner(rowValue)
	}
	return rSlice
}

func mapResult(r rows, returnType reflect.Type,
		cache func(reflect.Type,string)(func(reflect.Value)sql.Scanner, error)) ([]interface{}, error) {
	cols, err := r.Columns()
	if err != nil {
		return nil, err
	}
	scanners := make([]func(reflect.Value)sql.Scanner,len(cols))
	for i,col := range(cols){
		s, err := cache(returnType, col)
		scanners[i] = s
		if err != nil {
			return nil, err
		}
	}
	rowValues := make([]interface{},0,1)
	for r.Next() {
		rowValue := reflect.New(returnType)
		err := r.Scan(mapScanners(scanners, reflect.Indirect(rowValue))...)
		if err != nil {
			return nil, err
		}
		rowValues = append(rowValues,rowValue.Interface())
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return rowValues, nil
}
