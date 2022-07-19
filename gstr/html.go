package gstr

import "regexp"

var (
	reHtmlTag = regexp.MustCompile(`<[^<]+?>|\xa0`)
)

func RemoveHtmlTag(s string) (rs string) {
	return reHtmlTag.ReplaceAllString(s, "")
}
