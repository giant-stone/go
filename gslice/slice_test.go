package gslice

import (
	"reflect"
	"sort"
	"testing"

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

		{map[string]struct{}{"甲": {}, "乙": {}, "丙": {}}, []string{"甲", "乙", "丙"}},
		{map[string]struct{}{"foo": {}, "bar": {}}, []string{"foo", "bar"}},
		{map[string]struct{}{}, []string{}},

		{nil, []string{}},
	} {
		got := UniqMapToSlice(item.s)
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
		got := MergeSliceInUniq(item.a, item.b)
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
		got := SliceIndex(len(item.haystack), func(i int) bool {
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
		got := MergeSliceInUniqAndOrder(item.a, item.b)
		require.Equal(t, item.want, got)
	}
}

func TestUniqMapToSliceInt64(t *testing.T) {
	for _, item := range []struct {
		s    map[int64]struct{}
		want []int64
	}{
		{
			map[int64]struct{}{
				1: {},
				5: {},
				3: {},
			},
			[]int64{1, 3, 5},
		},

		{map[int64]struct{}{}, []int64{}},

		{nil, []int64{}},
	} {
		got := UniqMapToSliceInt64(item.s)
		sort.Sort(Int64Slice(got))
		sort.Sort(Int64Slice(item.want))
		require.Equal(t, item.want, got)
	}
}

func TestInsert(t *testing.T) {
	// Define the type for the test cases
	type args[T any] struct {
		slice   []T
		index   int
		newItem T
	}

	// Test with int type
	t.Run("int tests", func(t *testing.T) {
		intTests := []struct {
			name string
			args args[int]
			want []int
		}{
			{
				name: "Insert into int slice",
				args: args[int]{slice: []int{1, 2, 3, 4, 5}, index: 2, newItem: 10},
				want: []int{1, 2, 10, 3, 4, 5},
			},
			{
				name: "Insert at the beginning of int slice",
				args: args[int]{slice: []int{2, 3, 4}, index: 0, newItem: 1},
				want: []int{1, 2, 3, 4},
			},
			{
				name: "Insert at the end of int slice",
				args: args[int]{slice: []int{1, 2, 3}, index: 3, newItem: 4},
				want: []int{1, 2, 3, 4},
			},
			{
				name: "Insert at an out-of-bounds index",
				args: args[int]{slice: []int{1, 2, 3}, index: 5, newItem: 10},
				want: []int{1, 2, 3}, // It should not modify the slice
			},
		}

		// Run each int test case
		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Insert(tt.args.slice, tt.args.index, tt.args.newItem)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Insert() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test with string type
	t.Run("string tests", func(t *testing.T) {
		stringTests := []struct {
			name string
			args args[string]
			want []string
		}{
			{
				name: "Insert into string slice",
				args: args[string]{slice: []string{"a", "b", "c", "d"}, index: 1, newItem: "x"},
				want: []string{"a", "x", "b", "c", "d"},
			},
			{
				name: "Insert at the beginning of string slice",
				args: args[string]{slice: []string{"b", "c", "d"}, index: 0, newItem: "a"},
				want: []string{"a", "b", "c", "d"},
			},
			{
				name: "Insert at the end of string slice",
				args: args[string]{slice: []string{"a", "b", "c"}, index: 3, newItem: "d"},
				want: []string{"a", "b", "c", "d"},
			},
			{
				name: "Insert at an out-of-bounds index",
				args: args[string]{slice: []string{"a", "b", "c"}, index: 5, newItem: "z"},
				want: []string{"a", "b", "c"}, // It should not modify the slice
			},
		}

		// Run each string test case
		for _, tt := range stringTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Insert(tt.args.slice, tt.args.index, tt.args.newItem)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Insert() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test with float64 type
	t.Run("float64 tests", func(t *testing.T) {
		float64Tests := []struct {
			name string
			args args[float64]
			want []float64
		}{
			{
				name: "Insert into float64 slice",
				args: args[float64]{slice: []float64{1.1, 2.2, 3.3, 4.4}, index: 3, newItem: 9.9},
				want: []float64{1.1, 2.2, 3.3, 9.9, 4.4},
			},
			{
				name: "Insert at the beginning of float64 slice",
				args: args[float64]{slice: []float64{2.2, 3.3}, index: 0, newItem: 1.1},
				want: []float64{1.1, 2.2, 3.3},
			},
			{
				name: "Insert at the end of float64 slice",
				args: args[float64]{slice: []float64{1.1, 2.2, 3.3}, index: 3, newItem: 4.4},
				want: []float64{1.1, 2.2, 3.3, 4.4},
			},
			{
				name: "Insert at an out-of-bounds index",
				args: args[float64]{slice: []float64{1.1, 2.2, 3.3}, index: 5, newItem: 10.0},
				want: []float64{1.1, 2.2, 3.3}, // It should not modify the slice
			},
		}

		// Run each float64 test case
		for _, tt := range float64Tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Insert(tt.args.slice, tt.args.index, tt.args.newItem)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Insert() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
