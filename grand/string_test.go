package grand_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/giant-stone/go/grand"
)

func TestString(t *testing.T) {
	cases := []struct {
		len int
	}{
		{16},
		{32},
	}

	for _, item := range cases {
		got1 := grand.String(item.len)
		got2 := grand.String(item.len)

		require.Equal(t, item.len, len(got1))
		require.Equal(t, item.len, len(got2))
		require.True(t, got1 != got2)
	}

	got1 := grand.String(5)
	got2 := grand.String(5)
	require.NotEqual(t, got1, got2)
}
