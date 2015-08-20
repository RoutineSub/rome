package godapper

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type example struct {
	Id int64,
	Val string
}

func ExampleUse(){
	database, err := sql.Open("postgres",)
	db := Wrap(nil)
	db.Execute("create table example (id BigInt NOT NULL, val VARCHAR(255) NOT NULL)", map[string]interface{}{})
	db.Execute("insert into example (id,val) values({id},{val})", map[string]interface{}{
			"id":100,
			"val":"Test"
		})
	tx, _ := db.Begin()
	vals, _ := tx.Query("select id as Id, val as Val from example", map[string]interface{}{}, new(example))
	tx.Commit()
	for _,v := range(vals){
		e := v.(example)
		fmt.Println(e.Id, e.Val)
	}
}