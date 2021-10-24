package main

import (
	"fmt"
	"gopkg.in/go-playground/pool.v3"
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

type jobResult struct {
	JobIdx int
	RetStr string
}

func main() {

	arr := generateJobArr(jobNum)

	p := pool.NewLimited(uint(poolSize))
	defer p.Close()

	res := make([]string, len(arr))

	batch := p.Batch()
	go func() {
		for idx, s := range arr {
			jobIdx, jobStr := idx, s // Copy Value to avoid copy pointer in Submit function!
			batch.Queue(func(wu pool.WorkUnit) (interface{}, error) {
				if wu.IsCancelled() {
					// return values not used
					return nil, nil
				}
				return job(jobStr, jobIdx)
			})
		}
		// DO NOT FORGET THIS OR GOROUTINES WILL DEADLOCK
		// if calling Cancel() it calls QueueComplete() internally
		batch.QueueComplete()
	}()

	for jobResultWrapper := range batch.Results() {
		if err := jobResultWrapper.Error(); err != nil {
			panic(err)
		}

		jobResVal := jobResultWrapper.Value()
		result := jobResVal.(*jobResult)
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
func job(str string, jobIdx int) (*jobResult, error) {

	fmt.Printf("str: %s, JobIdx: %d\n", str, jobIdx)
	retStr := prefix + str

	//time.Sleep(time.Millisecond * 500) // 睡眠500ms，模拟耗时

	// Goroutine Number Check：
	fmt.Printf("index: %d, goroutine Num: %d\n", jobIdx, runtime.NumGoroutine())

	return &jobResult{RetStr: retStr, JobIdx: jobIdx}, nil
}

// 初始化测试数据
func generateJobArr(jobNum int) []string {
	arr := make([]string, 0)
	for i := 1; i < jobNum+1; i++ {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
