package main

import (
	"fmt"
	"runtime"
	"sync"
)

const (
	defaultStackSize = 4096
)

func callPanic() {
	panic("test panic")
}

// getCurrentGoroutineStack 获取当前Goroutine的调用栈，便于排查panic异常
func getCurrentGoroutineStack() string {
	var buf [defaultStackSize]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func task(arr *[]int, i int, wg *sync.WaitGroup, lock *sync.Mutex) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[panic] err: %v\nstack: %s\n", err, getCurrentGoroutineStack())
		}
		wg.Done()
	}()

	if i == 500 {
		callPanic()
	}

	lock.Lock()
	defer lock.Unlock()
	*arr = append(*arr, i)
}

func main() {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	arr := make([]int, 0)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go task(&arr, i, &wg, &lock)
	}
	wg.Wait()

	fmt.Println(len(arr))
}
