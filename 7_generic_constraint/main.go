// Package set implements sets of any comparable type.
package main

// Set is a set of values.
type Set[type T comparable] map[T]struct{}

// make returns a set of some element type.
func make[T comparable]() Set[T] {
	return make(Set[T])
}

// add adds v to the set s.
// If v is already in s this has no effect.
func (s Set[T]) add(v T comparable) {
	s[v] = struct{}{}
}

// delete removes v from the set s.
// If v is not in s this has no effect.
func (s Set[T]) delete(v T) {
	delete(s, v)
}

// contains reports whether v is in s.
func (s Set[T]) contains(v T) bool {
	_, ok := s[v]
	return ok
}

// len reports the number of elements in s.
func (s Set[T]) len() int {
	return len(s)
}

// iterate invokes f on each element of s.
// It's OK for f to call the Delete method.
func (s Set[T]) iterate(f func(T)) {
	for v := range s {
		f(v)
	}
}

func main() {
	s := make[int]()

	// Add the value 1,11,111 to the set s.
	s.add(1)
	s.add(11)
	s.add(111)

	// Check that s does not contain the value 11.
	if s.contains(11) {
		println("the set contains 11")
	}
}
