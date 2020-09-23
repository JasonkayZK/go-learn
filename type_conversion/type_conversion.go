package main

import "fmt"

type Printer interface {
	Print()
}

type Student struct {
	Name string
	Age int
}

func (s *Student) Print() {
	fmt.Println("hello")
}

func main() {
	nilPointer := (*Student)(nil)
	nilPointer.Print()
}
