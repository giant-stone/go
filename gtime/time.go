package gtime

import (
	"errors"
	"time"
)

var (
	ErrInvalidTime = errors.New("invalid time")
)

// Yyyymmdd2unixTimeUtc convert a string in format "2006-01-02" to unix time.
func Yyyymmdd2unixTimeUtc(s string) (rs int64, err error) {
	if s == "" {
		err = ErrInvalidTime
		return
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return
	}

	t, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return
	}

	rs = t.Unix()
	return
}

// UnixTime2YyyymmddUtc convert a unix time to a string in format "2006-01-02".
func UnixTime2YyyymmddUtc(t int64) (rs string) {
	return time.Unix(t, 0).UTC().Format("2006-01-02")
}
