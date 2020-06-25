package main

import (
	"fmt"
	"prime/prime"
	"time"
)

func main() {
	t1 := time.Now().UnixNano()

	prime.SerialPrintPrimeN(50000)
	t2 := time.Now().UnixNano()
	fmt.Printf("serial time: %d ns\n", t2 - t1)

	prime.PrintPrimeN(50000)
	t3 := time.Now().UnixNano()
	fmt.Printf("parallel time: %d ns\n", t3 - t2)

	prime.NewPrintPrimeN(50000)
	fmt.Printf("new parallel time: %d ns\n", time.Now().UnixNano() - t3)
}
