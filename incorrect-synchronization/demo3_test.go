package incorrect_synchronization

import (
	"testing"
)

var a3 string
var done3 bool

func setup3() {
	a3 = "hello, world"
	done3 = true
}

func TestDemo3(t *testing.T) {
	go setup3()
	for !done3 {
	}
	print(a3)
}
