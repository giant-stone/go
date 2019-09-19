package util

import (
	"log"
	"runtime/debug"
)

// CheckErr print error with stack context and return true for error else false.
func CheckErr(err error) bool {
	if err != nil {
		log.Println("[error]", err, string(debug.Stack()))
		return true
	}
	return false
}
