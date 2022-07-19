package gslice

import "strings"

// UniqMapToSlice converts a map into a slice of string.
func UniqMapToSlice(m map[string]struct{}) (rs []string) {
	rs = make([]string, 0)
	for item := range m {
		if item == "" {
			continue
		}
		rs = append(rs, item)
	}
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
