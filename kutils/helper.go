package kutils

import "encoding/json"

func ConvertSliceToString(s []interface{}) string {
	r, _ := json.Marshal(s)
	return string(r)
}
