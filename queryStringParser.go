package godapper

import (
	"bytes"
	"errors"
)

func parseQueryString(db database, queryString string) (*parsedQueryString, error) {
	var output bytes.Buffer
	var key bytes.Buffer
	//builders for getting expressions out of input, default to 2 so we don't get to big
	keyList := make([]string, 0, 2)
	writeMode := 0;
	for _,c := range queryString {
		switch writeMode {
		case 0: //regular write through
			switch c {
			case '{': //begin a replacement
				writeMode = 1
			default:
				output.WriteRune(c)
			}
		case 1: //replacement mode
			if c == '{' {
				//output a { character
				output.WriteRune('{')
				writeMode = 3 //close or error
				break //don't fall through
			}
			//Not a escaped '{' start to parse an identfier
			writeMode = 2
			fallthrough
		case 2: //identifier must start with a valid character or index
			switch c {
			case '}':
				output.WriteRune('?')
				keyList = append(keyList, key.String())
				key = bytes.Buffer{}
				writeMode = 0
			default:
				key.WriteRune(c)
			}
		case 3:
			switch c {
			case '}': //Done without Error
				writeMode = 0
			default:
				return nil,errors.New("No closing brace in brace literal") //Error
			}
		}
	}
	//create a prepared statement
	stmt, err := db.Prepare(output.String())
	if err != nil {
		return nil, err
	}
	return parsed(stmt, keyList), nil
}