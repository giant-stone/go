package gruntime_test

import (
	"testing"

	"github.com/giant-stone/go/gruntime"
	"github.com/stretchr/testify/require"
)

func myFunc() {}

func TestGetFunctionName(t *testing.T) {
	got := gruntime.GetFunctionName(myFunc)
	want := "github.com/giant-stone/go/gruntime_test.myFunc"
	require.Equal(t, want, got, "")
}
