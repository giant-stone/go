package gslice_test

import (
	"sort"
	"testing"

	"github.com/giant-stone/go/gslice"

	"github.com/stretchr/testify/require"
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
