// All Unix time(stamp) in UTC by default.
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

// UnixTimeToYMDHmUtc convert a unix time to a string in format "2006-01-02" in UTC.
func UnixTimeToYMDHmUtc(t int64) (rs string) {
	return time.Unix(t, 0).UTC().Format("2006-01-02")
}

// UnixTimeToYMDHmLocal convert a unix time to a string in format "2006-01-02" in specified timezone.
func UnixTimeToYMDHmLocal(t int64, tz *time.Location) (rs string) {
	return time.Unix(t, 0).In(tz).Format("2006-01-02")
}

// UnixTimeToYMDHmSUtc convert a unix time to a string in UTC in format "2006-01-02 15:04:05".
func UnixTimeToYMDHmSUtc(t int64) (rs string) {
	return time.Unix(t, 0).UTC().Format("2006-01-02 15:04:05")
}

// UnixTimeToYMDHmSLocal convert a unix time to a string in specified timezone in format "2006-01-02 15:04:05".
func UnixTimeToYMDHmSLocal(t int64, tz *time.Location) (rs string) {
	return time.Unix(t, 0).In(tz).Format("2006-01-02 15:04:05")
}

func MustParseDateInUnixTimeUtc(s string) (rs int64) {
	rs, _ = Yyyymmdd2unixTimeUtc(s)
	return rs
}

// UnixTime2YyyymmddUtc convert a unix time to a string in format "2006-01-02" in UTC.
// Deprecated: As of v1.1.0, use UnixTimeToYMDHmUtc instead.
func UnixTime2YyyymmddUtc(t int64) (rs string) {
	return UnixTimeToYMDHmUtc(t)
}

// UnixTime2YyyymmddLocal convert a unix time to a string in format "2006-01-02" in specified timezone.
// Deprecated: As of v1.1.0, use UnixTimeToYMDHmLocal instead.
func UnixTime2YyyymmddLocal(t int64, tz *time.Location) (rs string) {
	return UnixTimeToYMDHmLocal(t, tz)
}

// UnixTime2YYYYMMDDHHmmUtc convert a unix time to a string in UTC in format "2006-01-02 15:04:05".
// Deprecated: As of v1.1.0, use UnixTimeToYMDHmSUtc instead.
func UnixTime2YYYYMMDDHHmmUtc(t int64) (rs string) {
	return UnixTimeToYMDHmSUtc(t)
}

// UnixTime2YYYYMMDDHHmmLocal convert a unix time to a string in specified timezone in format "2006-01-02 15:04:05".
// Deprecated: As of v1.1.0, use UnixTimeToYMDHmSLocal instead.
func UnixTime2YYYYMMDDHHmmLocal(t int64, tz *time.Location) (rs string) {
	return UnixTimeToYMDHmSLocal(t, tz)
}

// Deprecated: As of v1.1.0, use MustParseDateInUnixTimeUtc instead.
func MustParseDateInUnixtimeUtc(s string) (rs int64) {
	return MustParseDateInUnixTimeUtc(s)
}
