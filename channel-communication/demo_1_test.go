package main

import "testing"

var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0
	//close(c)
}

func TestDemo1(t *testing.T) {
	go f()
	<-c
	print(a)
}
