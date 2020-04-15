package main

import (
    "database/sql"
    "fmt"
    "github.com/astaxie/beedb"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:123456@/test?charset=utf8")
    if err != nil {
        panic(err)
    }
    orm := beedb.New(db)

    var user User
    err = orm.Where("id=?", 1).Find(&user)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%v\n", user)
}

type User struct {
    Id       int `beedb:"PK"`
    Username string
    Password string
    Age      int
}
