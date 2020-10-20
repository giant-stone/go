package gutil

import (
	"encoding/json"
)

// Struct2map converts object from a struct into a map for use in gsql.
func Struct2map(obj interface{}) *map[string]interface{} {
	b, _ := json.Marshal(obj)
	m := map[string]interface{}{}
	_ = json.Unmarshal(b, &m)
	return &m
}
