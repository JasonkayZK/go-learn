package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var p *struct{} = nil
	fmt.Printf("%p\n", p)  //0x0
	fmt.Println(unsafe.Sizeof(p)) // 8

	var s []int = nil
	fmt.Printf("%p\n", s)  //0x0
	fmt.Println(unsafe.Sizeof(s)) // 24

	var m map[int]bool = nil
	fmt.Printf("%p\n", m)  //0x0
	fmt.Println(unsafe.Sizeof(m)) // 8

	var c chan string = nil
	fmt.Printf("%p\n", c)  //0x0
	fmt.Println(unsafe.Sizeof(c)) // 8

	var f func() = nil
	fmt.Printf("%p\n", f)  //0x0
	fmt.Println(unsafe.Sizeof(f)) // 8

	var i interface{} = nil
	fmt.Printf("%p\n", i)  //%!p(<nil>)
	fmt.Println(unsafe.Sizeof(i)) // 16
}
