package main

import "fmt"

func main() {
	// 123
	nil := 123
	fmt.Println(nil)
	//cannot use nil (type int) as type map[string]int in assignment
	//var _ map[string]int = nil
}
