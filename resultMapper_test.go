package dapper

import (
	"database/sql"
	"testing"
)

type mockRows struct {
	err func() error
	cols []string
	rows [][]interface{}
	index int
}

func (m mockRows) Err() error {
	return m.err()
}

func (m mockRows) Columns() ([]string, error) {
	return m.cols, nil
}

func (m *mockRows) Next() bool {
	m.index++;
	if m.index < len(m.rows) {
		return true;
	}
	return false
}

func (m *mockRows) Scan(dest ...interface{}) error {
	for i,col := range(m.rows[m.index]) {
		scanner, ok := dest[i].(sql.Scanner)
		if ok {
			err := scanner.Scan(col)
			if err != nil {
				return err
			}
		} else {
			dest[i] = col
		}
	}
	return nil
}

type test1 struct {
	A, B int
}

func TestMapResultReturnsListOfStructs(t *testing.T) {
	mr := mockRows{func()error{
			return nil
		}, []string{"A","B"},[][]interface{}{[]interface{}{1,2}}, -1}
	result := MapResult(&mr)
	dest := test1{}
	pos := 0
	for {
		err := result.Decode(&dest)
		if err != nil {
			break
		}
		if dest.A != 1 {
			t.Fatalf("Expected A value of 1 but was %v",dest.A)
		}
		if dest.B != 2 {
			t.Fatalf("Expected B value of 2 but was %v",dest.B)
		}
		dest = test1{}
		pos++
	}
	if pos != 1 {
		t.Fatalf("Expected 2 rows but was %v",pos)
	}
}

func TestMapResultWorksWithMaps(t *testing.T) {
	mr := mockRows{func()error {
		return nil
		}, []string{"A","B"},
		[][]interface{}{[]interface{}{"U","V"}},-1}
	result := MapResult(&mr)
	dest := make(map[string]interface{})
	len := 0
	for {
		err := result.Decode(dest)
		if err != nil {
			break
		}
		if dest["A"] != "U" {
			t.Fatalf("Expected A value of U but was %v",dest["A"])
		}
		if dest["B"] != "V" {
			t.Fatalf("Expected B value of V but was %v",dest["B"])
		}
		dest = make(map[string]interface{})
		len++
	}
	if len != 1 {
		t.Fatalf("Expected length of 1 but was %v",len)
	}
}

func TestMapResultWorksWithSlices(t *testing.T) {
	mr := mockRows{func()error{return nil},
			[]string{"A","B"},
			[][]interface{}{[]interface{}{1,"U"}}, -1}
	result := MapResult(&mr)
	dest := make([]interface{},2)
	len := 0
	for {
		err := result.Decode(dest)
		if err != nil {
			break
		}
		if dest[0] != 1 {
			t.Fatalf("Expected first column value of 1 but was %v", dest[0])
		}
		if dest[1] != "U" {
			t.Fatalf("Expected second column value of U but was %v", dest[1])
		}
		dest = make([]interface{},2)
		len++
	}
	if len != 1 {
		t.Fatalf("Expected length of 1 but was %v", len)
	}
}

func TestMapResultReturnsListOfStructsPerRow(t *testing.T) {
	mr := mockRows{func()error {
			return nil
		}, []string{"A","B"},
		[][]interface{}{
			[]interface{}{1,2},
			[]interface{}{3,4},
			[]interface{}{5,6},
			[]interface{}{7,8}}, -1}
	result := MapResult(&mr)
	dest := &test1{}
	pos := 0
	for {
		err := result.Decode(dest)
		if err != nil {
			break
		}
		if pos == 0 {
			if dest.A != 1 {
				t.Fatalf("Expected A value of 1 but was %v",dest.A)
			}
			if dest.B != 2 {
				t.Fatalf("Expected B value of 2 but was %v", dest.B)
			}
		} else if pos == 1 {
			if dest.A != 3{
				t.Fatalf("Expected A value of 3 but was %v", dest.A)
			}
			if dest.B != 4 {
				t.Fatalf("Expected B value of 4 but was %v", dest.B)
			}
		} else if pos == 2 {
			if dest.A != 5 {
				t.Fatalf("Expected A value of 5 but was %v",dest.A)
			}
			if dest.B != 6 {
				t.Fatalf("Expected B value of 6 but was %v", dest.B)
			}
		} else if pos == 3 {
			if dest.A != 7 {
				t.Fatalf("Expected A value of 7 but was %v",dest.A)
			}
			if dest.B != 8 {
				t.Fatalf("Expected B value of 8 but was %v",dest.B)
			}
		}
		dest = &test1{}
		pos++
	}
	if pos != 4 {
		t.Fatalf("Expected length of 4 but was %v",pos)
	}
}