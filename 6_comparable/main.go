package main

import (
	"fmt"
)

// index returns the index of x in s, or -1 if not found.
func index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable constraint
		// so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

type Foo struct {
	a string
	b int
}

func main() {
	fmt.Println(index([]int{1, 2, 3, 4, 5}, 3))
	fmt.Println(index([]string{"a", "b", "c", "d", "e"}, "d"))
	fmt.Println(index(
		[]Foo{
			{"a", 1},
			{"b", 2},
			{"c", 3},
			{"d", 4},
			{"e", 5},
		}, Foo{"b", 2}))
}
