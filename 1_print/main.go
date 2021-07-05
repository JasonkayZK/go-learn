package main

import (
	"fmt"
)

func printSlice[T any](s []T) {
	for _, v := range s {
		fmt.Printf("%v ", v)
	}
	fmt.Print("\n")
}

func main() {
	printSlice[int]([]int{1, 2, 3, 4, 5})
	printSlice[float64]([]float64{1.01, 2.02, 3.03, 4.04, 5.05})
	printSlice[string]([]string{"one", "two", "three", "four", "five"})

	printSlice([]int64{5, 4, 3, 2, 1})
}
