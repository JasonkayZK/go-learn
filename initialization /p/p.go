package p

import "fmt"

func init() {
	func() {
		fmt.Println("p init in a new goroutine")
	}()
}

func P() {
	fmt.Println("this is p")
}
