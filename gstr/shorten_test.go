package gstr_test

import (
	"reflect"
	"testing"

	"github.com/giant-stone/go/gstr"
)

const sampleShortenMaxLen = 5

var (
	samplesShorten = []struct {
		s    string
		want string
	}{
		{"foo", "foo"},
		{"", ""},
		{"helloworld", "hello"},
		{"秦始皇统一中国童男童女炼丹", "秦始皇统一"},
	}
)

func TestShorten(t *testing.T) {
	for _, item := range samplesShorten {
		got := gstr.Shorten(item.s, sampleShortenMaxLen)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Shorten -%s- want %s got %s", item.s, item.want, got)
		}
	}
}

var (
	samplesTrimSubstrings = []struct {
		s      string
		needle []string
		want   string
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
		got := gstr.TrimSubstrings(item.s, item.needle)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("TrimSubstrings -%s- want %s got %s", item.s, item.want, got)
		}
	}
}
