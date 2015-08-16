package godapper

import (
	"reflect"
	"testing"
)

type test1 struct {
	A, B int
}

func TestMapResultReturnsListOfStructs(t *testing.T) {
	mr := mockRows{func()error{
			return nil
		}, -1,[]string{"A","B"},[][]interface{}{[]interface{}{1,2}}}
	results,_ := mapResult(&mr,reflect.TypeOf(test1{}),getScanner)
	if len(results) != 1 {
		t.Fatalf("Expected 1 result but had %v", len(results))
	}
	r := results[0]
	v, ok := r.(*test1)
	if !ok {
		t.Fatalf("Expected *test1 type but was %v",reflect.TypeOf(r))
	}
	if v.A != 1 {
		t.Fatalf("Expected A value of 1 but was %v", v.A)
	}
	if v.B != 2 {
		t.Fatalf("Excepted B value of 2 but was %v", v.B)
	}
}

func TestMapResultReturnsListOfStructsPerRow(t *testing.T) {
	mr := mockRows{func()error {
			return nil
		}, -1, []string{"A","B"},
		[][]interface{}{
			[]interface{}{1,2},
			[]interface{}{3,4},
			[]interface{}{5,6},
			[]interface{}{7,8}}}
	results,_ := mapResult(&mr, reflect.TypeOf(test1{}),getScanner)
	if len(results) != 4 {
		t.Fatalf("Expected 4 results but had %v", len(results))
	}

}