// All Unix time(stamp) in UTC.
// https://en.wikipedia.org/wiki/Unix_time
package gtime

import (
	"errors"
	"time"
)

var (
	ErrInvalidTime = errors.New("invalid time")
)

// Yyyymmdd2unixTimeUtc convert a string in UTC in format "2006-01-02" to unix time.
func Yyyymmdd2unixTimeUtc(s string) (rs int64, err error) {
	if s == "" {
		return 0, ErrInvalidTime
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return 0, err
	}

	t, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return 0, err
	}

	rs = t.Unix()
	return rs, nil
}

// UnixTime2YyyymmddUtc convert a unix time to a string in format "2006-01-02" in UTC.
func UnixTime2YyyymmddUtc(t int64) (rs string) {
	return time.Unix(t, 0).UTC().Format("2006-01-02")
}

// UnixTime2YyyymmddLocal convert a unix time to a string in format "2006-01-02" in specified timezone.
func UnixTime2YyyymmddLocal(t int64, tz *time.Location) (rs string) {
	return time.Unix(t, 0).In(tz).Format("2006-01-02")
}

// UnixTime2YYYYMMDDHHmmUtc convert a unix time to a string in UTC in format "2006-01-02 15:04:05".
func UnixTime2YYYYMMDDHHmmUtc(t int64) (rs string) {
	return time.Unix(t, 0).UTC().Format("2006-01-02 15:04:05")
}

// UnixTime2YYYYMMDDHHmmLocal convert a unix time to a string in specified timezone in format "2006-01-02 15:04:05".
func UnixTime2YYYYMMDDHHmmLocal(t int64, tz *time.Location) (rs string) {
	return time.Unix(t, 0).In(tz).Format("2006-01-02 15:04:05")
}

func MustParseDateInUnixtimeUtc(s string) (rs int64) {
	rs, _ = Yyyymmdd2unixTimeUtc(s)
	return rs
}
