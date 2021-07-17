package main

import (
	"fmt"
)

type addable interface {
	type int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64, complex64, complex128
}

func generator[T addable](a T, v T) func() T {
	return func() T {
		r := a
		a = a + v
		return r
	}
}

func main() {
	g1 := generator(0, 1)
	fmt.Println(g1())
	fmt.Println(g1())
	fmt.Println(g1())

	g2 := generator(-9.9, 0.1)
	fmt.Println(g2())
	fmt.Println(g2())
	fmt.Println(g2())
}
