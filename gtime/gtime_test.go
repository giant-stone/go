package gtime_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/giant-stone/go/gtime"
	"github.com/stretchr/testify/require"
)

func TestYyyymmdd2unixTimeUtc(t *testing.T) {
	samples := []struct {
		s       string
		want    int64
		wantErr error
	}{
		{"1999-08-07", 933984000, nil},
		{"", 0, gtime.ErrInvalidTime},
	}

	for _, tc := range samples {
		got, gotErr := gtime.Yyyymmdd2unixTimeUtc(tc.s)
		require.ErrorIs(t, gotErr, tc.wantErr, tc.s)
		require.Equal(t, tc.want, got, tc.s)
	}
}

func TestMustParseDateInUnixTimeUtc(t *testing.T) {
	samples := []struct {
		s    string
		want int64
	}{
		{"1999-08-07", 933984000},
		{"", 0},
	}

	for _, tc := range samples {
		got := gtime.MustParseDateInUnixTimeUtc(tc.s)
		require.Equal(t, tc.want, got, tc.s)
	}
}

func TestUnixTimeToYMDHmUtc(t *testing.T) {
	samples := []struct {
		s    int64
		want string
	}{
		// 933984001 <=> "1999-08-07 00:00:01" in UTC
		{933984001, "1999-08-07"},
		// 1673478001 <=> "2023-01-11 23:00:01" in UTC
		{1673478001, "2023-01-11"},
	}

	for _, tc := range samples {
		got := gtime.UnixTimeToYMDHmUtc(tc.s)

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("UnixTimeToYMDHmUtc -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}

func TestUnixTimeToYMDHmLocal(t *testing.T) {
	samples := []struct {
		s    int64
		tz   *time.Location
		want string
	}{
		// 1673478001 <=> "2023-01-11 23:00:01" in UTC
		{1673478001, time.FixedZone("UTC+8", 8*3600), "2023-01-12"},
		{1673478001, time.FixedZone("America/New_York", -5*3600), "2023-01-11"},
	}

	for _, tc := range samples {
		got := gtime.UnixTimeToYMDHmLocal(tc.s, tc.tz)

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("UnixTimeToYMDHmLocal -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}

func TestUnixTimeToYMDHmSUtc(t *testing.T) {
	samples := []struct {
		s    int64
		want string
	}{
		{1673478001, "2023-01-11 23:00:01"},
	}

	for _, tc := range samples {
		got := gtime.UnixTimeToYMDHmSUtc(tc.s)

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("UnixTimeToYMDHmSUtc -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}

func TestUnixTimeToYMDHmSLocal(t *testing.T) {
	samples := []struct {
		s    int64
		tz   *time.Location
		want string
	}{
		// 1673478001 <=> "2023-01-11 23:00:01" in UTC
		{1673478001, time.FixedZone("UTC+8", 8*3600), "2023-01-12 07:00:01"},
		{1673478001, time.FixedZone("America/New_York", -5*3600), "2023-01-11 18:00:01"},
	}

	for _, tc := range samples {
		got := gtime.UnixTimeToYMDHmSLocal(tc.s, tc.tz)

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("UnixTimeToYMDHmSLocal -%v- want %v got %v", tc.s, tc.want, got)
		}
	}
}
