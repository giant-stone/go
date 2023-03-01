package gstrconv

import (
	"regexp"
	"strconv"
)

// Atoi parses a string into int.
func Atoi(s string) (i int) {
	v, _ := strconv.Atoi(s)
	return v
}

// Atoi64 parses a string into int64.
func Atoi64(s string) (rs int64) {
	i, _ := strconv.ParseUint(s, 10, 64)
	return int64(i)
}

var (
	reDigit = regexp.MustCompile(`(?P<digit>\d+)`)
)

// ParseDigitFromMixed parses a digit from a string contains 0-9 and non 0-9 chars.
func ParseDigitFromMixed(s string) (i int) {
	return Atoi(reDigit.FindString(s))
}

// ParseFloat32 parses a string into float32.
func ParseFloat32(s string) (i float32) {
	v, _ := strconv.ParseFloat(s, 32)
	return float32(v)
}
