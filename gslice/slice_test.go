package gslice_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/giant-stone/go/gslice"
)

func TestUniqMapToSlice(t *testing.T) {
	for _, item := range []struct {
		s    map[string]struct{}
		want []string
	}{
		{
			map[string]struct{}{
				"a": {},
				"c": {},
				"b": {},
			},
			[]string{"a", "b", "c"},
		},

		{map[string]struct{}{"甲": {}, "乙": {}, "丙": {}}, []string{"甲", "乙", "丙"}},
		{map[string]struct{}{"foo": {}, "bar": {}}, []string{"foo", "bar"}},
		{map[string]struct{}{}, []string{}},

		{nil, []string{}},
	} {
		got := gslice.UniqMapToSlice(item.s)
		sort.Strings(got)
		sort.Strings(item.want)
		require.Equal(t, item.want, got)
	}
}

func TestMergeSliceInUniq(t *testing.T) {
	for _, item := range []struct {
		a    []string
		b    []string
		want map[string]struct{}
	}{
		{
			[]string{"a", "c"},
			[]string{"b", "c", "d"},
			map[string]struct{}{
				"a": {},
				"b": {},
				"c": {},
				"d": {},
			},
		},

		{nil, nil, map[string]struct{}{}},
	} {
		got := gslice.MergeSliceInUniq(item.a, item.b)
		require.Equal(t, item.want, got)
	}
}

func TestSliceIndex(t *testing.T) {
	for _, item := range []struct {
		haystack []string
		needle   string
		want     int
	}{
		{[]string{"foo", "bar", "baz"}, "go", -1},
		{[]string{"中文输入法", "输入法", "中文"}, "中文", 2},
	} {
		got := gslice.SliceIndex(len(item.haystack), func(i int) bool {
			return item.haystack[i] == item.needle
		})
		require.Equal(t, item.want, got, item.haystack)
	}
}

func TestMergeSliceInUniqAndOrder(t *testing.T) {
	for _, item := range []struct {
		a    []string
		b    []string
		want []string
	}{
		{
			[]string{"c", "a"},
			[]string{"b", "c", "d"},
			[]string{"c", "a", "b", "d"},
		},

		{
			[]string{"c", "b"},
			[]string{"c", "b"},
			[]string{"c", "b"},
		},

		{
			[]string{"a"},
			[]string{},
			[]string{"a"},
		},

		{
			[]string{},
			[]string{"b"},
			[]string{"b"},
		},

		{nil, nil, []string{}},
	} {
		got := gslice.MergeSliceInUniqAndOrder(item.a, item.b)
		require.Equal(t, item.want, got)
	}
}
