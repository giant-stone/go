package gstr_test

import (
	"testing"

	"github.com/giant-stone/go/gstr"
	"github.com/giant-stone/go/gutil"
)

var (
	samplesMd5 = []struct {
		s        string
		expected string
	}{
		{
			"hello",
			"5d41402abc4b2a76b9719d911017c592",
		},

		{
			"中文",
			"a7bac2239fcdcb3a067903d8077c4a07",
		},
	}
)

func TestMd5(t *testing.T) {
	for _, item := range samplesMd5 {
		gutil.CmpExpectedGot(t, item.s, item.expected, gstr.Md5(item.s))
	}
}
