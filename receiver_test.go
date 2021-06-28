package main

import (
	"fmt"
	"testing"
)

type Object struct {
	Elem string
}

func (o *Object) refEqual(o2 *Object) bool {
	fmt.Printf("refEqual str: %v\n", o.Elem == o2.Elem)
	return o == o2
}

func (o Object) copyEqual(o2 *Object) bool {
	fmt.Printf("copyEqual str: %v\n", o.Elem == o2.Elem)
	return &o == o2
}

func TestEqual(t *testing.T) {
	elem := Object{Elem: "haha"}

	fmt.Printf("ref func: %v\n", elem.refEqual(&elem))
	fmt.Printf("copy func: %v\n", elem.copyEqual(&elem))
}
