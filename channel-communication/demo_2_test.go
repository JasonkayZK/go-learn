package main

import "testing"

var c2 = make(chan int, 10)
var a2 string

func f2() {
	a2 = "hello, world"
	<-c2
}

func TestDemo2(t *testing.T) {
	go f2()
	c2 <- 0
	print(a2)
}
