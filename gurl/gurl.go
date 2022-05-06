package gurl

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func Params2qs(params map[string][]interface{}) (rs string) {
	keysInOrder := make(sort.StringSlice, 0)
	for key := range params {
		keysInOrder = append(keysInOrder, key)
	}

	sort.Sort(keysInOrder)

	pairs := make([]string, 0)
	for _, key := range keysInOrder {
		values := params[key]

		value := strings.Join(sliceInterface2sliceStr(values), ",")
		pair := fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
		pairs = append(pairs, pair)
	}

	return strings.Join(pairs, "&")
}

func sliceInterface2sliceStr(items []interface{}) (rs sort.StringSlice) {
	rs = make(sort.StringSlice, 0)
	for _, item := range items {
		rs = append(rs, fmt.Sprintf("%v", item))
	}
	sort.Sort(rs)
	return
}

func ParseQueryStringByName(s, name string) (rs string) {
	var qs string
	splits := strings.Split(s, "?")
	if len(splits) > 1 {
		qs = splits[1]
	} else {
		qs = s
	}

	if !strings.Contains(qs, "&") {
		rs = s
	} else {
		rsParse, _ := url.ParseQuery(qs)
		rs = rsParse.Get(name)
	}
	return
}
