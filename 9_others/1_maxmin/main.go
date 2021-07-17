package main

import "fmt"

type comparableNum interface {
	type int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		uintptr, float32, float64
}

func max[T comparableNum](a []T) T {
	m := a[0]
	for _, v := range a {
		if m < v {
			m = v
		}
	}
	return m
}

func min[T comparableNum](a []T) T {
	m := a[0]
	for _, v := range a {
		if m > v {
			m = v
		}
	}
	return m
}

func main() {
	v1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	maxRes1 := max[int](v1)
	minRes1 := min[int](v1)
	fmt.Printf("v1: max: %d, min: %d\n", maxRes1, minRes1)

	v2 := []float64{1.1, 2.3, 3.1, 4.4, 5.5, 6.2, 7.1, 8.9, 9.2, 10.23}
	maxRes2 := max[float64](v2)
	minRes2:= min[float64](v2)
	fmt.Printf("v2: max: %f, min: %f\n", maxRes2, minRes2)
}
