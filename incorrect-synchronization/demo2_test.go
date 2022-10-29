package incorrect_synchronization

import (
	"sync"
	"testing"
)

var a2 string
var done bool
var once sync.Once

func setup() {
	a2 = "hello, world"
	done = true
}

func doprint() {
	if !done {
		once.Do(setup)
	}
	print(a2)
}

func twoprint() {
	go doprint()
	go doprint()
}

func TestDemo2(t *testing.T) {
	twoprint()
}
