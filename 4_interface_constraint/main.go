package main

import (
	"fmt"
	"strconv"
)

type StringInt int

func (i StringInt) String() string {
	return strconv.Itoa(int(i))
}

type MyStringer interface {
	String() string
}

// A generic type that need to implement MyStringer interface!
func stringify[T MyStringer](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.String())
	}
	return ret
}

func main() {
	fmt.Println(stringify([]StringInt{1, 2, 3, 4, 5}))
}
