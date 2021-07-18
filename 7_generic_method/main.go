// Package set implements sets of any comparable type.
package main

import "fmt"

type addable interface {
	//type int, int8, int16, int32, int64,
	//	uint, uint8, uint16, uint32, uint64, uintptr,
	//	float32, float64, complex64, complex128,
	comparable
}

// set is a set of values.
type set[T addable] map[T]struct{}

// add adds v to the set s.
// If v is already in s this has no effect.
func (s set[T]) add(v T) {
	s[v] = struct{}{}
}

// contains reports whether v is in s.
func (s set[T]) contains(v T) bool {
	_, ok := s[v]
	return ok
}

// len reports the number of elements in s.
func (s set[T]) len() int {
	return len(s)
}

// delete removes v from the set s.
// If v is not in s this has no effect.
func (s set[T]) delete(v T) {
	delete(s, v)
}

// iterate invokes f on each element of s.
// It's OK for f to call the Delete method.
func (s set[T]) iterate(f func (T) ) {
	for v := range s {
		f(v)
	}
}

// invalid AST: method must have no type parameters
// methods cannot have type parameters
/*
func (s set[T]) anotherGeneric[P comparable](v T, p P) {
	fmt.Printf("v: %v, p: %v\n", v, p)
}
*/

func print[T addable](s T) {
	fmt.Printf("%v ", s)
}

func main() {
	s := make(set[int])

	// Add the value 1,11,111 to the set s.
	s.add(1)
	s.add(11)
	s.add(111)
	s.add(1111)
	s.add(11111)
	fmt.Printf("%v\n", s)

	// Check that s does not contain the value 11.
	if s.contains(11) {
		println("the set contains 11")
	} else {
		println("the set do not contain 11")
	}

	// Check len of set
	fmt.Printf("the len of set: %d\n", s.len())

	// Delete elem in set
	s.delete(11)
	fmt.Println("\nafter delete 11:")
	if s.contains(11) {
		println("the set contains 11")
	} else {
		println("the set do not contain 11")
	}
	fmt.Printf("the len of set: %d\n", s.len())

	// Iterate set with explicit type(int)
	s.iterate(func(x int) {
		fmt.Println(x + 1)
	})

	// Iterate set with implicit type(addable) ERROR!
	//s.iterate(print)

	// Generic in another generic type method
	//s.anotherGeneric(2, 3)
}
