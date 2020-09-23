package main

import (
	"fmt"
)

func main() {
	fmt.Println((map[string]int)(nil) == nil) //true
	fmt.Println((func())(nil) == nil)         //true
}
