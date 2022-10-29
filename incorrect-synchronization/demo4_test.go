package incorrect_synchronization

import (
	"testing"
)

type T struct {
	msg string
}

var g4 *T

func setup4() {
	t := new(T)
	t.msg = "hello, world"
	g4 = t
}

func TestDemo4(t *testing.T) {
	go setup4()
	for g4 == nil {
	}
	print(g4.msg)
}
