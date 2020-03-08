package gref_test

import (
	"log"
	"os"
	"testing"

	"github.com/giant-stone/go/gref"
)

func myFunc() {}

func TestGetFunctionName(t *testing.T) {
	name := gref.GetFunctionName(myFunc)
	nameExpected := "github.com/giant-stone/go/gref_test.myFunc"
	if name != nameExpected {
		t.Errorf("expected name=myFunc, got %s", nameExpected)
	}
}


func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	os.Exit(m.Run())
}
