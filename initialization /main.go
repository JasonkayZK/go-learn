package main

import (
	_ "initialization/q"

	_ "initialization/p"
)

// Will print:
//  p init in a new goroutine
//  q init in a new goroutine
// Even though the q is imported first!
func main() {
}
