package gunicode

import "unicode"

// IsInteger check a string if is an integer.
func IsInteger(s string) (rs bool) {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsDigit(rune(r)) {
			return false
		}
	}
	return true
}
