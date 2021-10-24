package main

import (
	"fmt"
	"github.com/panjf2000/ants"
	"sync"
)

type jobItem struct {
	Str    string
	JobIdx int
	Res    *[]string
	Wg     *sync.WaitGroup
}

func antsWithFunc() {

	arr := generateJobArr(jobNum)

	funcPool, err := ants.NewPoolWithFunc(poolSize,
		func(i interface{}) {
			item := i.(*jobItem)
			job(item.Str, item.JobIdx, item.Res, item.Wg)
		}, func(opts *ants.Options) {
			opts.Nonblocking = false
			opts.MaxBlockingTasks = len(arr)
		})
	if err != nil {
		panic(err)
	}
	defer funcPool.Release()

	wg := sync.WaitGroup{}
	res := make([]string, len(arr))

	for idx, s := range arr {
		jobIdx, jobStr := idx, s // Copy Value to avoid copy pointer in Submit function!
		wg.Add(1)
		err := funcPool.Invoke(&jobItem{
			Str:    jobStr,
			JobIdx: jobIdx,
			Res:    &res,
			Wg:     &wg,
		})
		if err != nil {
			panic(fmt.Errorf("submit job err: %v", err))
		}
	}
	wg.Wait()

	// Result Test
	for idx, re := range res {
		if re != prefix+arr[idx] {
			panic(fmt.Sprintf("not equal: re: %s, arr[%d]: %s", re, idx, arr[idx]))
		}
	}
}
