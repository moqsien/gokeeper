package kutils

import (
	"encoding/json"

	"github.com/gogf/gf/util/gconv"
)

// SliceToString 切片转字符串
func SliceToString[T string | interface{}](s []T) string {
	var t T
	var typ interface{} = t
	switch typ.(type) {
	case string:
		r := ""
		for i, v := range s {
			if i != 0 {
				r += "," + gconv.String(v)
			} else {
				r += gconv.String(v)
			}
		}
		return r
	default:
		r, _ := json.Marshal(s)
		return string(r)
	}
}
