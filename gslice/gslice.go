package gslice

// MapToSlice convert a map to a slice, map in format { string => struct{} }.
func MapToSlice(m *map[string]struct{}) (rs []string) {
	all := make([]string, 0)
	for item := range *m {
		all = append(all, item)
	}
	return all
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
