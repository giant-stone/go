package gtime_test

import (
	"reflect"
	"testing"

	"github.com/giant-stone/go/gtime"
	"github.com/stretchr/testify/require"
)

func TestYyyymmdd2unixTimeUtc(t *testing.T) {
	var (
		samples = []struct {
			s       string
			want    int64
			wantErr error
		}{
			{"1999-08-07", 933984000, nil},
			{"", 0, gtime.ErrInvalidTime},
		}
	)

	for _, item := range samples {
		got, gotErr := gtime.Yyyymmdd2unixTimeUtc(item.s)
		require.ErrorIs(t, gotErr, item.wantErr, item.s)
		require.Equal(t, item.want, got, item.s)
	}
}

func TestMustParseDateInUnixtimeUtc(t *testing.T) {
	var (
		samples = []struct {
			s    string
			want int64
		}{
			{"1999-08-07", 933984000},
			{"", 0},
		}
	)

	for _, item := range samples {
		got := gtime.MustParseDateInUnixtimeUtc(item.s)
		require.Equal(t, item.want, got, item.s)
	}
}

func TestUnixTime2YyyymmddUtc(t *testing.T) {
	var (
		samples = []struct {
			s    int64
			want string
		}{
			{933984000, "1999-08-07"},
		}
	)

	for _, item := range samples {
		got := gtime.UnixTime2YyyymmddUtc(item.s)

		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("UnixTime2YyyymmddUtc -%v- want %v got %v", item.s, item.want, got)
		}
	}
}