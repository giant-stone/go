package gslice

import (
	"sort"
	"strings"
)

// UniqMapToSlice converts a map into a slice of string.
func UniqMapToSlice(m map[string]struct{}) (rs []string) {
	rs = make([]string, 0)
	for item := range m {
		if item == "" {
			continue
		}
		rs = append(rs, item)
	}
	sort.StringSlice(rs).Sort()
	return rs
}

// MergeSliceInUniq merges two slice of string into map in unique.
func MergeSliceInUniq(a []string, b []string) (rs map[string]struct{}) {
	rs = make(map[string]struct{}, 0)
	for _, item := range a {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		rs[item] = struct{}{}
	}
	for _, item := range b {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		rs[item] = struct{}{}
	}
	return
}

// SliceIndex get item index of a slice
//  ref:https://stackoverflow.com/a/8307594/913751
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
