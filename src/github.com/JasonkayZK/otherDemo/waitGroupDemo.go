package main

import (
	"fmt"
	"sync"
)

func main() {
	s := "sss"

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		s = "aaa"
		fmt.Printf("Output in Closure: %s\n", s)
		waitGroup.Done()
	}()


	//waitGroup.Wait() // s = sss
	fmt.Printf("Output in main: %s\n", s)
	//waitGroup.Wait() // s = aaa
}
