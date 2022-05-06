package gruntime

import (
	"reflect"
	"runtime"
)

// GetFunctionName returns function full name.
// https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
