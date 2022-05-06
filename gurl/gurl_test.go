package gurl_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/giant-stone/go/gurl"
)

func TestParams2qs(t *testing.T) {
	samples := []struct {
		params map[string][]interface{}
		want   string
	}{
		{
			map[string][]interface{}{
				"name": {"foo"},
				"age":  {123},
			},
			"age=123&name=foo",
		},

		{
			map[string][]interface{}{
				"name": {"张三"},
				"dept": {"研发"},
			},
			"dept=" + url.QueryEscape("研发") + "&name=" + url.QueryEscape("张三"),
		},
	}
	for _, item := range samples {
		got := gurl.Params2qs(item.params)
		require.Equal(t, item.want, got, item.params)
	}
}

func TestParseQueryStringByName(t *testing.T) {
	samples := []struct {
		s    string
		name string
		want string
	}{
		{"foo.com/?name=foo&age=123", "name", "foo"},
		{"?name=foo&age=123", "name", "foo"},
		{"name=foo&age=123", "name", "foo"},
		{"name=foo&age=123", "age", "123"},
	}
	for _, item := range samples {
		got := gurl.ParseQueryStringByName(item.s, item.name)
		require.Equal(t, item.want, got, item.s)
	}
}
