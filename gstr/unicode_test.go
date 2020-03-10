package gstr_test

import (
	"log"
	"os"
	"testing"

	"github.com/giant-stone/go/gstr"
	"github.com/giant-stone/go/gutil"
)

var (
	samplesIsInteger = []struct {
		S        string
		Expected bool
	}{
		{"123", true},
		{"a", false},
		{" 123 ", true},
		{"a 123", false},
		{"中文", false},
		{"123中文", false},
	}
)

func TestIsInteger(t *testing.T) {
	for _, sample := range samplesIsInteger {
		gutil.CmpExpectedGot(t, sample.S, sample.Expected, gstr.IsInteger(sample.S))
	}

}

var (
	samplesIsAscii = []struct {
		S        rune
		Expected bool
	}{
		{'1', true},
		{'a', true},
		{' ', false},
		{'⌘', false},
		{'中', false},
	}
)

func TestIsAscii(t *testing.T) {
	for _, sample := range samplesIsAscii {
		gutil.CmpExpectedGot(t, sample.S, sample.Expected, gstr.IsAscii(sample.S))
	}
}

var (
	samplesIsAlphabet = []struct {
		S        rune
		Expected bool
	}{
		{'1', false},
		{'a', true},
		{' ', false},
		{'⌘', false},
		{'中', false},
	}
)

func TestIsAlphabet(t *testing.T) {
	for _, sample := range samplesIsAlphabet {
		gutil.CmpExpectedGot(t, sample.S, sample.Expected, gstr.IsAlphabet(sample.S))
	}
}

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	os.Exit(m.Run())
}
