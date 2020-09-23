package main

import (
	"fmt"
	"time"
)

var stop chan interface{}

func reqTask(name string) {
	for {
		select {
		case <-stop:
			fmt.Println("stop", name)
			return
		default:
			fmt.Println(name, "send request")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	stop = make(chan interface{})
	go reqTask("worker1")
	time.Sleep(3 * time.Second)
	stop <- struct{}{}
	time.Sleep(3 * time.Second)
}
