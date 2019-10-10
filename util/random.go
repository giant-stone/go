package gutil

import (
	"math/rand"
	"time"
)

func RandChoice(objs []interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	return objs[rand.Intn(len(objs))]
}
