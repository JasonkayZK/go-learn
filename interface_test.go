package main

import (
	"fmt"
	"testing"
)

type Introduce interface {
	Introduce() string
}

type Item struct {
	Elem string
}

func (i *Item) Introduce() string {
	return "haha from reference"
}

//func (i Item) Introduce() string {
//	return "haha from object"
//}

func TestInterface(t *testing.T) {
	elem := &Item{Elem: "haha"}

	fmt.Println(elem.Introduce())
	fmt.Println((*elem).Introduce())
}
