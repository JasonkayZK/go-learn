package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type Foo struct {
	a int
}

func main() {
	for i := 0; i < 3; i++ {
		f := NewFoo(i)
		println(f.a)
	}

	runtime.GC()

	time.Sleep(time.Second)
}

//go:noinline
func NewFoo(i int) *Foo {
	f := &Foo{a: rand.Intn(50)}
	runtime.SetFinalizer(f, func(f *Foo) {
		fmt.Println(`foo ` + strconv.Itoa(i) + ` has been garbage collected`)
	})

	return f
}
