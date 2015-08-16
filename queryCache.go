package godapper

import (
	"reflect"
	"database/sql"
)

type mapperKey struct {
	t reflect.Type
	columnName string
}

type cache struct {
	queries map[string]*parsedQueryString
	mappers map[mapperKey]func(reflect.Value)sql.Scanner
	parse func(database,string)(*parsedQueryString,error)
	build func(reflect.Type,string)(func(reflect.Value)sql.Scanner, error)
}

func (c *cache) getCachedQuery(db database, queryString string) (*parsedQueryString, error) {
	pqs, ok := c.queries[queryString]
	if !ok {
		pqs, err := c.parse(db, queryString)
		if err != nil {
			return nil, err
		}
		c.queries[queryString] = pqs
	}
	return pqs, nil
}

func (c *cache) getCachedMapper(t reflect.Type, columnName string) (func(reflect.Value)sql.Scanner, error) {
	key := mapperKey{t,columnName}
	mps, ok := c.mappers[key]
	if !ok {
		mps, err := c.build(t, columnName)
		if err != nil {
			return nil, err
		}
		c.mappers[key] = mps
	}
	return mps, nil
}
