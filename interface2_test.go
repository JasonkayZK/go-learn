package main

import (
	"fmt"
	"testing"
)

type SelfCompare interface {
	Equal(item *MyElem) bool
}

type MyElem struct {
	Elem string
}

//func (m *MyElem) Equal(item *MyElem) bool {
//	return m == item
//}

func (m MyElem) Equal(item *MyElem) bool {
	return &m == item
}

func TestMyElem(t *testing.T) {
	elem := &MyElem{Elem: "haha"}

	fmt.Printf("ref func: %v\n", elem.Equal(elem))
}

func callEqual(s SelfCompare) {
	fmt.Println(s.Equal(s.(*MyElem)))
}

func TestInterfaceCall(t *testing.T) {
	elem := &MyElem{Elem: "haha"}
	fmt.Println(elem.Equal(elem))
	callEqual(elem)
}
