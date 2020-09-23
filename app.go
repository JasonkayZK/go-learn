package main

import (
	"fmt"
)

type Student struct {
	Name string
	Age  int
}

func (s *Student) Say() {
	s.sayHello()
}

func (s *Student) sayHello() {
	fmt.Println("hello")
}

func (s *Student) SayName() {
	fmt.Println(s.Name)
}

func main() {
	//student, err := InitStudent(true)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//student.Say()

	var nilPointer *Student = nil
	nilPointer.Say()
}

func InitStudent(needErr bool) (*Student, error) {
	if needErr {
		return nil, fmt.Errorf("init student err")
	}

	return &Student{}, nil
}
