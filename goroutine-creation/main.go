package main

import (
	"fmt"
	"time"
)

var a string

func f() {
	fmt.Println(a)
}

func hello() {
	a = "hello, world"
	go f()
}

func main() {

	hello()

	fmt.Println("hello returned")

	<-time.After(1 * time.Second)
}
