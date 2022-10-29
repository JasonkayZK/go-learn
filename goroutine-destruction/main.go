package main

import (
	"fmt"
	"time"
)

var a string

func hello() {
	go func() { a = "hello" }()
	fmt.Println(a)
}

func main() {

	hello()

	fmt.Println("hello returned")

	<-time.After(1 * time.Second)
}
