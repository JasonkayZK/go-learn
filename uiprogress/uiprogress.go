package main

import (
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	//simple()
	//customDecorators()
	multi()
}

func simple() {
	uiprogress.Start()            // start rendering
	bar := uiprogress.AddBar(100) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	bar.PrependElapsed()

	for bar.Incr() {
		time.Sleep(time.Millisecond * 20)
	}
}

func customDecorators() {
	uiprogress.Start() // start rendering

	var steps = []string{
		"downloading source",
		"installing deps",
		"compiling",
		"packaging",
		"seeding database",
		"deploying",
		"staring servers",
	}
	bar := uiprogress.AddBar(len(steps))

	// prepend the current step to the bar
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "app: " + steps[b.Current()-1]
	})

	for bar.Incr() {
		time.Sleep(time.Millisecond * 1000)
	}
}

func multi() {
	waitTime := time.Millisecond * 100
	uiprogress.Start()

	var wg sync.WaitGroup

	bar1 := uiprogress.AddBar(20).AppendCompleted().PrependElapsed()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar1.Incr() {
			time.Sleep(waitTime)
		}
	}()

	bar2 := uiprogress.AddBar(40).AppendCompleted().PrependElapsed()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar2.Incr() {
			time.Sleep(waitTime)
		}
	}()

	time.Sleep(time.Second)
	bar3 := uiprogress.AddBar(20).PrependElapsed().AppendCompleted()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar3.Incr() {
			time.Sleep(waitTime)
		}
	}()

	wg.Wait()
}
