package main

import (
	"sync"
	"testing"
)

var l sync.Mutex
var a string

func f() {
	a = "hello, world"
	l.Unlock()
}

func TestDemo1(t *testing.T) {
	l.Lock()
	go f()
	l.Lock()
	print(a)
}
