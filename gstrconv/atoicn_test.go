package gstrconv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/giant-stone/go/gstrconv"
)

func TestAtoiCn(t *testing.T) {
	for idx, tc := range []struct {
		s    string
		want int64
		err  error
	}{
		{"5", 5, nil},
		{"20", 20, nil},
		{"100010", 100_010, nil},
		{"999,999,999,999", 999_999_999_999, nil},
		{"999_999_999_999", 999_999_999_999, nil},

		{"一", 1, nil},
		{"十", 10, nil},
		{"十一", 11, nil},
		{"二十二", 22, nil},
		{"一百", 100, nil},
		{"一百零一", 101, nil},

		{"一千零五", 1005, nil},
		{"一亿三千万", 130_000_000, nil},

		{"一亿零三百万零一百五十一", 103_000_151, nil},
		{"一亿零三百零一万零一百五十一", 103_010_151, nil},

		{"老友记", 0, gstrconv.ErrInvalidData},
		{"foo", 0, gstrconv.ErrInvalidData},
		{"a", 0, gstrconv.ErrInvalidData},
		{"", 0, gstrconv.ErrInvalidData},
	} {
		got, gotErr := gstrconv.AtoiCn(tc.s)
		require.ErrorIs(t, tc.err, gotErr, "idx=%d", idx)
		if tc.err == nil {
			require.Equal(t, tc.want, got, "idx=%d s=-%s- want %d got %d", idx, tc.s, tc.want, got)
		}
	}
}
