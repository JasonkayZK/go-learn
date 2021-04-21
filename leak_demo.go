package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func query() int {
	n := rand.Intn(100)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return n
}

// 每次执行此函数，都会导致有两个goroutine处于阻塞状态
func queryAll() int {
	ch := make(chan int)
	go func() { ch <- query() }()
	go func() { ch <- query() }()
	go func() { ch <- query() }()
	// <-ch
	// <-ch
	return <-ch
}

func main() {
	// 每次循环都会泄漏两个goroutine
	for i := 0; i < 4; i++ {
		queryAll()
		// main()也是一个主goroutine
		fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
	}
}
