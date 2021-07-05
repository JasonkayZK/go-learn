package main

import (
	"fmt"
	"strconv"
)

// Compiling Error:
// StringInt does not satisfy MySignedStringer (uint not found in int, int8, int16, int32, int64)
//type StringInt uint

type StringInt int

func (i StringInt) String() string {
	return strconv.Itoa(int(i))
}

type MySignedStringer interface {
	// Only these types & interface can be generalized!
	type int, int8, int16, int32, int64
	String() string
}

// A generic type that need to implement MySignedStringer interface!
func stringify[T MySignedStringer](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.String())
	}
	return ret
}

func main() {
	fmt.Println(stringify([]StringInt{1, 2, 3, 4, 5}))
}
