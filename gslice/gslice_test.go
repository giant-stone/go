package gslice_test

import (
	"sort"
	"testing"

	"github.com/giant-stone/go/gslice"
	"github.com/stretchr/testify/require"
)

func TestMapToSlice(t *testing.T) {
	for _, item := range []struct {
		s    map[string]struct{}
		want []string
	}{
		{map[string]struct{}{"甲": {}, "乙": {}, "丙": {}}, []string{"甲", "乙", "丙"}},
		{map[string]struct{}{"foo": {}, "bar": {}}, []string{"foo", "bar"}},
		{map[string]struct{}{}, []string{}},
	} {
		got := gslice.MapToSlice(&item.s)
		sort.StringSlice(item.want).Sort()
		require.Equal(t, item.want, got, item.s)
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
