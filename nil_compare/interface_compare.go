package main

import (
	"fmt"
)

func main() {
	// false
	fmt.Println((interface{})(nil) == (*int)(nil))
}
