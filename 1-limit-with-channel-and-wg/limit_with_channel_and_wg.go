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
	jobNum = 1000000

	poolSize = runtime.NumCPU() << 1
)

type jobReqItem struct {
	Str    string
	JobIdx int
	Res    *[]string
	Wg     *sync.WaitGroup
}

func main() {

	arr := generateJobArr(jobNum)
	fmt.Printf("generate job arr finished: %d\n", len(arr))
	time.Sleep(time.Second * 3)

	startTime := time.Now().UnixNano() / 1000

	wg := sync.WaitGroup{}
	jobChan := make(chan *jobReqItem, poolSize)
	res := make([]string, len(arr))

	fmt.Printf("start job, pool size: %d\n", poolSize)

	// Start Consumer: 生成指定数目的 goroutine，每个 goroutine 消费 jobsChan 中的数据
	for i := 0; i < poolSize; i++ {
		go func() {
			for jobReq := range jobChan {
				job(jobReq)
			}
		}()
	}

	// Start Producer: 把 job 依次推送到 jobsChan 供 goroutine 消费
	for idx, s := range arr {
		wg.Add(1)
		jobChan <- &jobReqItem{Str: s, JobIdx: idx, Res: &res, Wg: &wg}
		fmt.Printf("index: %d, goroutine Num: %d \n", idx, runtime.NumGoroutine())
	}
	wg.Wait()
	close(jobChan)
	fmt.Printf("end job, goroutine Num: %d, cost time: %dms\n", runtime.NumGoroutine(), time.Now().UnixNano()/1000-startTime)

	// Test
	for idx, re := range res {
		if re != prefix+arr[idx] {
			panic(fmt.Sprintf("not equal: re: %s, arr[idx]: %s", re, arr[idx]))
		}
	}
}

// 任务内容
func job(jobReq *jobReqItem) {
	defer jobReq.Wg.Done()

	fmt.Printf("str: %s, jobIdx: %d\n", jobReq.Str, jobReq.JobIdx)
	(*jobReq.Res)[jobReq.JobIdx] = prefix + jobReq.Str

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
