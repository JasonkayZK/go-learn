package main

import (
	"fmt"
	"testing"
	"time"
)

var limit = make(chan int, 3)

func TestDemo3(t *testing.T) {

	work := make([]func(), 0)
	for i := 0; i < 10; i++ {
		aI := i
		work = append(work, func() {
			fmt.Printf("Hello, this is: %d\n", aI)
		})
	}

	for _, w := range work {
		go func(w func()) {
			limit <- 1
			w()
			time.Sleep(time.Second * 1)
			<-limit
		}(w)
	}

	<-time.After(time.Second * 10)
}
