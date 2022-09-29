package grand

import (
	"math/rand"
	"time"
)

// Rand returns range in [0~n)
func Rand(n int) int {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r1.Intn(n)
}
