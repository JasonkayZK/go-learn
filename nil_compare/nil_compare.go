package main

import "fmt"

func main() {
	nilCompare2()
}

func nilCompare() {
	//var m map[int]string
	//var ptr *int

	//invalid operation: m == ptr (mismatched types map[int]string and *int)
	//fmt.Printf(m == ptr)
}

func nilCompare2() {
	type IntPtr *int
	fmt.Println(IntPtr(nil) == (*int)(nil))        //true
	fmt.Println((interface{})(nil) == (*int)(nil)) //false
}
