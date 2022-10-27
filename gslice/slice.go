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

type Int64Slice []int64

func (x Int64Slice) Len() int           { return len(x) }
func (x Int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// UniqMapToSliceInt64 converts a map into a slice of int64.
func UniqMapToSliceInt64(m map[int64]struct{}) (rs []int64) {
	rs = make([]int64, 0)
	for item := range m {
		if item == 0 {
			continue
		}
		rs = append(rs, item)
	}
	sort.Sort(Int64Slice(rs))
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

// MergeSliceInUniqAndOrder merges two slice of string into map in unique.
func MergeSliceInUniqAndOrder(merged []string, part []string) (rs []string) {
	sorted := make([]string, 0)
	unique := make(map[string]struct{}, 0)

	for _, item := range merged {
		if _, dup := unique[item]; dup {
			continue
		}
		unique[item] = struct{}{}
		sorted = append(sorted, item)
	}

	for _, item := range part {
		if _, dup := unique[item]; dup {
			continue
		}
		unique[item] = struct{}{}
		sorted = append(sorted, item)
	}

	return sorted
}
