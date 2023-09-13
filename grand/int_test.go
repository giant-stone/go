package grand_test

import (
	"fmt"
	"testing"

	"github.com/giant-stone/go/grand"
)

func ExampleRand() {
	fmt.Println(grand.Rand(10))
	fmt.Println(grand.Rand(3))
}

func ExampleRandRange() {
	fmt.Println(grand.RandRange(1_000, 9_999))
	fmt.Println(grand.RandRange(100_0000, 999_9999))
}

func TestRandRange(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "1_000~9_999",
			args: args{
				min: 1_000,
				max: 9_000,
			},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := grand.RandRange(tt.args.min, tt.args.max); got < tt.args.min || got > tt.args.max {
				t.Errorf("RandRange() = %v, want %v < got < %v", got, tt.args.min, tt.args.max)
			}
		})
	}
}
