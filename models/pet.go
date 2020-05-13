package models

import (
	"time"
)

type (
	Pet struct {
		Id uint `json:"id"`
		Name string `json:"name" xorm:"varchar(20)"`
		Age int `json:"age"`
		Photo string `json:"photo" xorm:"varchar(30)"`
		CreatedAt time.Time `json:"created_at" xorm:"created"`
		UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	}
)
