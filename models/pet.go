package models

import (
	"time"
)

type (
	Pet struct {
		Id    int       `json:"id" binding:"required" form:"id"`
		Name  string    `json:"name" xorm:"varchar(20)" binding:"required" form:"name"`
		Age   int       `json:"age" binding:"required" form:"age"`
		Photo string    `json:"photo" xorm:"varchar(30)" form:"photo"`
		Ctime time.Time `json:"created_at" xorm:"ctime"`
	}
)
