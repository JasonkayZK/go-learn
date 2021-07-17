package main

import (
	"fmt"
)

func ternaryOperation[T any](s bool, right T, wrong T) T {
	if s {
		return right
	}
	return wrong
}

func main() {
	fmt.Println(ternaryOperation(4 < 5, "less", "greater"))
	fmt.Println(ternaryOperation(5 < 4, 5, 4))
}
