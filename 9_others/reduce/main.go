package main

import (
	"fmt"
)

// reduceFunc
// a: reduce operate collection
// f: reduce operation function
// initial: initial value in reduce
func reduceFunc[T any](a []T, f func(T, T) T, initial interface{}) T {
	if len(a) == 0 || f == nil {
		var vv T
		return vv
	}

	l := len(a) - 1
	reduce := func(a []T, ff func(T, T) T, memo T, startPoint, direction, length int) T {
		result := memo
		index := startPoint
		for i := 0; i <= length; i++ {
			result = ff(result, a[index])
			index += direction
		}
		return result
	}

	if initial == nil {
		return reduce(a, f, a[0], 1, 1, l-1)
	}

	return reduce(a, f, initial.(T), 0, 1, l)
}

func main() {
	v1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res1 := reduceFunc(v1, func(lhs, rhs int) int {
		return lhs + rhs
	}, 1)
	fmt.Println(res1)

	v2 := []string{"x", "y", "z"}
	res2 := reduceFunc(v2, func(lhs, rhs string) string {
		return lhs + rhs
	}, "a")
	fmt.Println(res2)

	v3 := []int{5, 4, 3, 2, 1}
	res3 := reduceFunc(v3, func(lhs, rhs int) int {
		return lhs*10 + rhs
	}, 0)
	fmt.Println(res3)
}
