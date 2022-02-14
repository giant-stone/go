package gstr_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/giant-stone/go/gstr"
)

var (
	samplesIsInteger = []struct {
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
)

func TestIsInteger(t *testing.T) {
	for _, item := range samplesIsInteger {
		got := gstr.IsInteger(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("IsInteger -%v- want %v got %v", item.s, item.want, got)
		}
	}
}

var (
	samplesIsAscii = []struct {
		s    rune
		want bool
	}{
		{'1', true},
		{'a', true},
		{' ', false},
		{'⌘', false},
		{'中', false},
	}
)

func TestIsAscii(t *testing.T) {
	for _, item := range samplesIsAscii {
		got := gstr.IsAscii(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("IsAscii -%v- want %v got %v", item.s, item.want, got)
		}
	}
}

var (
	samplesIsAlphabet = []struct {
		s    rune
		want bool
	}{
		{'1', false},
		{'a', true},
		{' ', false},
		{'⌘', false},
		{'中', false},
	}
)

func TestIsAlphabet(t *testing.T) {
	for _, item := range samplesIsAlphabet {
		got := gstr.IsAlphabet(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("IsAlphabet -%v- want %v got %v", item.s, item.want, got)
		}
	}
}

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	os.Exit(m.Run())
}
