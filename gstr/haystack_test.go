package gstr_test

import (
	"testing"

	"github.com/giant-stone/go/gstr"
	"github.com/giant-stone/go/gutil"
)

var (
	samplesStrInSlice = []struct {
		haystack []string
		needle   string
		expected bool
	}{
		{
			[]string{"a", "b", "c"},
			"a",
			true,
		},

		{
			[]string{"英国", "加拿大", "澳大利亚"},
			"英国",
			true,
		},

		{
			[]string{"usa", "加拿大", "澳大利亚"},
			"uk",
			false,
		},
	}
)

func TestStrInSlice(t *testing.T) {
	for _, sample := range samplesStrInSlice {
		gutil.CmpExpectedGot(t, "needle="+sample.needle, sample.expected, gstr.StrInSlice(sample.haystack, sample.needle))
	}

}
