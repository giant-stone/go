package gutil_test

import (
	"log"
	"testing"

	"github.com/giant-stone/go/gutil"
)

func TestStruct2map(t *testing.T) {
	type user struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	u := user{Id: 1234, Name: "foo"}
	m := gutil.Struct2map(&u)
	log.Println(u, m)

	// NOTE: interger auto convert into float64 in Go JSON
	idI, ok := (*m)["id"].(float64)
	idGot := int(idI)
	if !ok || idGot != u.Id {
		t.Errorf("want id=%v got %v", u.Id, idGot)
	}
	nameGot, ok := (*m)["name"].(string)
	if !ok || nameGot != u.Name {
		t.Errorf("want name=%v got %v", u.Name, nameGot)
	}
}
