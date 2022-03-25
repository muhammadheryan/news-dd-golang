package helper

import (
	"encoding/json"
)

func StructToJson(v interface{}) string {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(out)
}
