package main

import (
	"fmt"
	"go-restful-xorm/routing"
)

func main() {
	server := routing.WebService{}
	server.Run()
	fmt.Println("Server started!")
}
