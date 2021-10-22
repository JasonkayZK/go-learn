package main

import (
	"fmt"
	"runtime"
	"sync"
)

func antsWithFunc() {
	arr := generateJobArr(jobNum)

	wg := sync.WaitGroup{}
	res := make([]string, len(arr))

	for idx, s := range arr {
		wg.Add(1)
		go job(s, idx, &res, &wg)

		// Goroutine Number Check：
		// +1：包括了main函数的Goroutine
		// 两倍poolSize：是最差情况下，所有的Goroutine的锁全部释放的同时，所有新的Goroutine被创建
		fmt.Printf("index: %d, goroutine Num: %d\n", idx, runtime.NumGoroutine())
	}
	wg.Wait()

	// Result Test
	for idx, re := range res {
		if re != prefix+arr[idx] {
			panic(fmt.Sprintf("not equal: re: %s, arr[idx]: %s", re, arr[idx]))
		}
	}
}
