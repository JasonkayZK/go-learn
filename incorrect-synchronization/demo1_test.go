package incorrect_synchronization

import "testing"

var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func TestDemo1(t *testing.T) {
	go f()
	g()
}
