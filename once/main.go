package main

import (
	"fmt"
	"sync"
	"time"
)

var a string
var once sync.Once

func setup() {
	a = "hello, world"
	fmt.Println("a has been setup")
}

func doPrint() {
	once.Do(setup)
	fmt.Println(a)
}

func twoPrint() {
	go doPrint()
	go doPrint()
}

func main() {

	twoPrint()

	<-time.After(1 * time.Second)
}
