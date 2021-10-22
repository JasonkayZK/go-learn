package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"runtime"
	"strconv"
	"sync"
)

const (
	prefix = `hello: `
)

var (
	jobNum = 1000000

	// Fatal Err!
	//jobNum = 1000000000

	poolSize = runtime.NumCPU() // 同时运行的goroutine上限
)

func main() {

	arr := generateJobArr(jobNum)

	wg := sync.WaitGroup{}
	sem := semaphore.NewWeighted(int64(poolSize))
	res := make([]string, len(arr))

	for idx, s := range arr {
		err := sem.Acquire(context.Background(), 1)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go job(s, idx, &res, &wg, sem)

		// Goroutine Number Check：
		// +1：包括了main函数的Goroutine
		// 两倍poolSize：是最差情况下，所有的Goroutine的锁全部释放的同时，所有新的Goroutine被创建
		fmt.Printf("index: %d, goroutine Num: %d\n", idx, runtime.NumGoroutine())
		if runtime.NumGoroutine() > poolSize<<1+1 {
			panic("超过了指定Goroutine池大小！")
		}
	}
	wg.Wait()

	// Result Test
	for idx, re := range res {
		if re != prefix+arr[idx] {
			panic(fmt.Sprintf("not equal: re: %s, arr[idx]: %s", re, arr[idx]))
		}
	}
}

// 任务内容
func job(str string, jobIdx int, res *[]string, wg *sync.WaitGroup, sem *semaphore.Weighted) {
	defer func() {
		wg.Done()
		sem.Release(1)
	}()

	fmt.Printf("str: %s, jobIdx: %d\n", str, jobIdx)
	(*res)[jobIdx] = prefix + str

	//time.Sleep(time.Millisecond * 500) // 睡眠500ms，模拟耗时
}

// 初始化测试数据
func generateJobArr(jobNum int) []string {
	arr := make([]string, 0)
	for i := 1; i < jobNum+1; i++ {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
