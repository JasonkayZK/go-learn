package main

import (
	"fmt"
	"go-init-map/model"
	"go-init-map/model/repo"
)

func main() {
	repo.Students["test4"] = model.Student{Id: 4, Name: "test4", Grade: 4}

	for name, student := range repo.Students {
		fmt.Printf("Student name: %s, student: %v\n", name, student)
	}
}
