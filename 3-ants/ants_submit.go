package main

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants"
)

func antsSubmit() {

	arr := generateJobArr(jobNum)

	pool, err := ants.NewPool(poolSize, func(opts *ants.Options) {
		opts.Nonblocking = false
		opts.MaxBlockingTasks = len(arr)
	})
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	wg := sync.WaitGroup{}
	res := make([]string, len(arr))

	for idx, s := range arr {
		err := pool.Submit(func() {
			wg.Add(1)
			job(s, idx, &res, &wg)
		})
		if err != nil {
			panic(fmt.Errorf("submit job err: %v", err))
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
