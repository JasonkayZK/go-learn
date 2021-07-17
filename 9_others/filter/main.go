package main

import (
	"fmt"
)

func filterFunc[T any](a []T, f func(T) bool) []T {
	var n []T
	for _, e := range a {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func main() {
	v1 := []int{1, 2, 3, 4, 5, 6}
	v1 = filterFunc(v1, func(v int) bool {
		return v < 4.0
	})
	fmt.Println(v1)

	v2 := []float64{2.1, 3.2, 23.2, 2.3}
	v2 = filterFunc(v2, func(v float64) bool {
		return v < 4.0
	})
	fmt.Println(v2)
}
