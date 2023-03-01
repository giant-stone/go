package gunicode_test

import (
	"reflect"
	"testing"

	"github.com/giant-stone/go/gunicode"
)

func TestIsInteger(t *testing.T) {
	samples := []struct {
		s    string
		want bool
	}{
		{"123", true},

		{"a", false},
		{" 123 ", false},
		{"a 123", false},
		{"中文", false},
		{"123中文", false},
	}

	for _, tc := range samples {
		got := gunicode.IsInteger(tc.s)
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("IsInteger -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}
