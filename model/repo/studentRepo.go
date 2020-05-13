package repo

import (
	m "go-init-map/model"
)

var Students map[string]m.Student

func init() {
	Students = make(map[string]m.Student)
	Students["test1"] = m.Student{Id: 1, Name: "test1", Grade: 1}
	Students["test2"] = m.Student{Id: 2, Name: "test2", Grade: 2}
	Students["test3"] = m.Student{Id: 3, Name: "test3", Grade: 3}
}
