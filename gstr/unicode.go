package gstr

import (
	"strconv"
	"strings"
)

// IsInteger tests a string if it is a number \d+.
func IsInteger(s string) bool {
	value, err := strconv.Atoi(strings.TrimSpace(s))
	return err == nil && value > 0
}

// IsAscii tests a rune if it is a ascii character [a-zA-Z0-9].
func IsAscii(r rune) bool {
	return IsAlphabet(r) || (r >= '0' && r <= '9')
}

// IsAlphabet tests a rune if it is a character in alphabet [a-zA-Z].
func IsAlphabet(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
