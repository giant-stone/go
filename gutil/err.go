package gutil

import (
	"runtime/debug"

	"github.com/giant-stone/go/glogging"
)

// CheckErr print error with stack context and return true for error else false.
func CheckErr(err error) bool {
	if err != nil {
		glogging.Sugared.Error(err, " ", string(debug.Stack()))
		return true
	}
	return false
}

// ExitOnErr print fatal error with stack context and exit.
func ExitOnErr(err error) {
	if err != nil {
		glogging.Sugared.Fatal(err, " ", string(debug.Stack()))
	}
}
