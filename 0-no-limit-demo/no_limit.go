package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	prefix = `hello: `
)

var (
	jobNum = 100000

	// Fatal Err!
	//jobNum = 1000000000
)

func main() {

	arr := generateJobArr(jobNum)

	wg := sync.WaitGroup{}
	res := make([]string, len(arr))

	for idx, s := range arr {
		wg.Add(1)
		go job(s, idx, &res, &wg)
		fmt.Printf("index: %d, goroutine Num: %d \n", idx, runtime.NumGoroutine())
	}
	wg.Wait()

	for idx, re := range res {
		if re != prefix+arr[idx] {
			panic(fmt.Sprintf("not equal: re: %s, arr[%d]: %s", re, idx, arr[idx]))
		}
	}
}

// 任务内容
func job(str string, jobIdx int, res *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("str: %s, jobIdx: %d\n", str, jobIdx)
	(*res)[jobIdx] = prefix + str

	time.Sleep(time.Second * 5) // 睡眠5s，模拟耗时
}

// 初始化测试数据
func generateJobArr(jobNum int) []string {
	arr := make([]string, 0)
	for i := 1; i < jobNum+1; i++ {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
