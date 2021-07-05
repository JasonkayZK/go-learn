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

type vector[T any] []T

func main() {
	// Compiling Error
	// Cannot use generic type vector[T interface{}] without instantiation
	//vs0 := vector{1,2,3,4,5}

	vs := vector[int]{5,4,2,1}
	printSlice(vs)

	vs2 := vector[string]{"haha", "hehe"}
	printSlice(vs2)
}
