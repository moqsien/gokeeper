package kutils

import "encoding/json"

func ConvertSliceToString(s []string) string {
	r, _ := json.Marshal(s)
	return string(r)
}
