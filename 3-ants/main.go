package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

const (
	prefix = `hello: `
)

var (
	jobNum = 100000000

	poolSize = runtime.NumCPU()
)

func main() {
	//antsSubmit()

	antsWithFunc()
}

// 任务内容
func job(str string, jobIdx int, res *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("str: %s, jobIdx: %d\n", str, jobIdx)
	(*res)[jobIdx] = prefix + str
	//fmt.Printf("res: %v\n", *res)

	//time.Sleep(time.Millisecond * 500) // 睡眠500ms，模拟耗时

	// Goroutine Number Check：
	fmt.Printf("index: %d, goroutine Num: %d\n", jobIdx, runtime.NumGoroutine())
}

// 初始化测试数据
func generateJobArr(jobNum int) []string {
	arr := make([]string, 0)
	for i := 1; i < jobNum+1; i++ {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
