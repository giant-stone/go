package grand_test

import (
	"fmt"

	"github.com/giant-stone/go/grand"
)

func ExampleRand() {
	fmt.Println(grand.Rand(10))
	fmt.Println(grand.Rand(3))
}
