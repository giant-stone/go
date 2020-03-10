package gstr

import (
	"strings"
)

// Shorten it cuts and concatenates a part of string with ellipsis if its length long than maxLen in unicode.
func Shorten(b []byte, maxLen int) string {
	var bodyChunk []rune
	body := []rune(string(b))
	if len(body) > maxLen {
		bodyChunk = body[:maxLen]
		bodyChunk = append(bodyChunk, []rune("...")...)
	} else {
		bodyChunk = body
	}
	return string(bodyChunk)
}

// TrimSubstrings remove all needles item from haystack.
func TrimSubstrings(haystack string, needles []string) (result string) {
	for _, item := range needles {
		haystack = strings.ReplaceAll(haystack, item, "")
	}
	return haystack
}
