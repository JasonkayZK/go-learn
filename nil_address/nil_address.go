package main

import (
	"fmt"
)

func main() {
	var m map[int]string
	var ptr *int
	var sl []int
	var i interface{} = nil
	fmt.Printf("%p\n", m)   //0x0
	fmt.Printf("%p\n", ptr) //0x0
	fmt.Printf("%p\n", sl)  //0x0
	fmt.Printf("%p\n", i)   //%!p(<nil>)
}
