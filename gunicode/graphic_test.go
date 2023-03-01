package gunicode_test

import (
	"reflect"
	"testing"

	"github.com/giant-stone/go/gunicode"
)

func TestIsWordCharacters(t *testing.T) {
	samples := []struct {
		s    rune
		want bool
	}{
		{'1', true},
		{'a', true},

		{' ', false},
		{'⌘', false},
		{'中', false},
	}

	for _, tc := range samples {
		got := gunicode.IsWordCharacters(tc.s)
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("IsAscii -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}

func TestIsAlphabet(t *testing.T) {
	samples := []struct {
		s    rune
		want bool
	}{
		{'a', true},

		{'1', false},
		{' ', false},
		{'⌘', false},
		{'中', false},
	}

	for _, tc := range samples {
		got := gunicode.IsAlphabet(tc.s)
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("IsAlphabet -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}
