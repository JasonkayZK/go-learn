package main

import (
	"fmt"
	"math/rand"
	"time"
)

func shuffle[T any](a []T) {
	n := len(a)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	v1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffle(v1)
	fmt.Println(v1)

	v2 := []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"}
	shuffle(v2)
	fmt.Println(v2)

	shuffle(v2)
	fmt.Println(v2)
}
