package q

import (
	"fmt"
	"initialization/p"
)

func init() {
	func() {
		fmt.Println("q init in a new goroutine")
	}()
}

func RunP() {
	p.P()
}
