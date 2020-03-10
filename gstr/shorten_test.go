package gstr_test

import (
	"testing"

	"github.com/giant-stone/go/gstr"
	"github.com/giant-stone/go/gutil"
)

const sampleShortenMaxLen = 5

var (
	samplesShorten = []struct {
		s        string
		expected string
	}{
		{"foo", "foo"},
		{"", ""},
		{"helloworld", "hello..."},
		{"秦始皇统一中国童男童女炼丹", "秦始皇统一..."},
	}
)

func TestShorten(t *testing.T) {
	for _, item := range samplesShorten {
		gutil.CmpExpectedGot(t, "", item.expected, gstr.Shorten([]byte(item.s), sampleShortenMaxLen))
	}
}

var (
	samplesTrimSubstrings = []struct {
		haystack string
		needle   []string
		expected string
	}{
		{
			"language zh-cn zh-tw zh-hk",
			[]string{"zh-cn", "zh-hk", "en-us"},
			"language  zh-tw ",
		},
	}
)

func TestTrimSubstrings(t *testing.T) {
	for _, item := range samplesTrimSubstrings {
		gutil.CmpExpectedGot(t, "", item.expected, gstr.TrimSubstrings(item.haystack, item.needle))
	}
}
