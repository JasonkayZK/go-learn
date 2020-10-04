package slice

import (
	"fmt"
	"testing"
)

func TestSlice1(t *testing.T) {
	array := [6]int64{1, 2, 3, 4, 5, 6}
	sliceA := array[2:5:5]
	sliceB := array[1:3:5]

	fmt.Println(sliceA, len(sliceA), cap(sliceA))
	fmt.Println(sliceB, len(sliceB), cap(sliceB))
}

