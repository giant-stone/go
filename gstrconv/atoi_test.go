package gstrconv_test

import (
	"reflect"
	"testing"

	"github.com/giant-stone/go/gstrconv"
)

func TestAtoi(t *testing.T) {
	for _, item := range []struct {
		s    string
		want int
	}{
		{
			"2008",
			2008,
		},

		{
			"2008-1-1",
			2008,
		},

		{
			"2008-12-30",
			2008,
		},
	} {
		got := gstrconv.ParseDigitFromMixed(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Shorten -%v- want %v got %v", item.s, item.want, got)
		}
	}

}

func TestParseDigitFromMixed(t *testing.T) {
	for _, item := range []struct {
		s    string
		want int
	}{
		{` S2
		E1
	`, 2},

		{` Episode 1`, 1},
		{"", 0},
	} {
		got := gstrconv.ParseDigitFromMixed(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Shorten -%v- want %v got %v", item.s, item.want, got)
		}
	}
}
