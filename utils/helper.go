package utils

import (
	"encoding/json"
	"fmt"
)

func StructToString(v interface{}) string {
	value, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return ""
	}
	return string(value)
}

func LogStruct(v interface{}) {
	fmt.Println(StructToString(v))
}
