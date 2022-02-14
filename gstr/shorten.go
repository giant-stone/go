package gstr

import (
	"strings"
	"unicode/utf8"
)

// TrimSubstrings remove all needles item from haystack.
func TrimSubstrings(haystack string, needles []string) (result string) {
	for _, item := range needles {
		haystack = strings.ReplaceAll(haystack, item, "")
	}
	return haystack
}

const (
	DefaultShortenSuffix = "…"
)

// Shorten it cuts and concatenates a part of string with ellipsis if its length long than maxLen in unicode.
func Shorten(s string, n int) (rs string) {
	strLen := utf8.RuneCountInString(s)

	if strLen <= n {
		rs = s
	} else {
		rs = Substring(s, 0, n)
	}
	return
}

func LenUtf8(s string) (n int) {
	return utf8.RuneCountInString(s)
}

func ShortenWith(s string, n int, suffix string) (rs string) {
	strLen := utf8.RuneCountInString(s)

	if strLen <= n {
		rs = s
	} else {
		rs = Substring(s, 0, n-1) + suffix
	}
	return
}

// https://stackoverflow.com/questions/28718682/how-to-get-a-substring-from-a-string-of-runes-in-golang
func Substring(s string, start int, end int) string {
	start_str_idx := 0
	i := 0
	for j := range s {
		if i == start {
			start_str_idx = j
		}
		if i == end {
			return s[start_str_idx:j]
		}
		i++
	}
	return s[start_str_idx:]
}
