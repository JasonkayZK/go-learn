package main

import (
	"fmt"
	"github.com/Jeffail/tunny"
	"runtime"
	"strconv"
)

const (
	prefix = `hello: `
)

var (
	jobNum = 100000

	poolSize = runtime.NumCPU()
)

type jobItem struct {
	Str    string
	JobIdx int
}

type jobResult struct {
	JobIdx int
	RetStr string
	Err    error
}

func main() {

	arr := generateJobArr(jobNum)

	pool := tunny.NewFunc(poolSize, func(jobItemEntity interface{}) interface{} {
		item := jobItemEntity.(*jobItem)
		return job(item.Str, item.JobIdx)
	})
	defer pool.Close()

	res := make([]string, len(arr))
	for idx, s := range arr {

		// Funnel this work into our pool. This call is synchronous and will
		// block until the job is completed.
		result := pool.Process(&jobItem{
			Str:    s,
			JobIdx: idx,
		}).(*jobResult)
		if result.Err != nil {
			panic(result.Err)
		}

		res[result.JobIdx] = result.RetStr
	}

	// Result Test
	for idx, re := range res {
		if re != prefix+arr[idx] {
			panic(fmt.Sprintf("not equal: re: %s, arr[%d]: %s", re, idx, arr[idx]))
		}
	}
}

// 任务内容
func job(str string, jobIdx int) *jobResult {

	fmt.Printf("str: %s, jobIdx: %d\n", str, jobIdx)
	retStr := prefix + str

	//time.Sleep(time.Millisecond * 500) // 睡眠500ms，模拟耗时

	// Goroutine Number Check：
	// +1：包括了main函数的Goroutine
	fmt.Printf("index: %d, goroutine Num: %d\n", jobIdx, runtime.NumGoroutine())
	if runtime.NumGoroutine() > poolSize+1 {
		panic("超过了指定Goroutine池大小！")
	}

	return &jobResult{RetStr: retStr, JobIdx: jobIdx, Err: nil}
}

// 初始化测试数据
func generateJobArr(jobNum int) []string {
	arr := make([]string, 0)
	for i := 1; i < jobNum+1; i++ {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
