package main

import (
	"fmt"
	"time"

	sqle "github.com/dolthub/go-mysql-server"

	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

const (
	user = "user"
	passwd = "pass"
	address = "localhost"
	port = "13306"
	dbName    = "test"
	tableName = "my_table"
)

func main() {
	engine := sqle.NewDefault()
	engine.AddDatabase(createTestDatabase())

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%s", address, port),
		Auth:     auth.NewNativeSingle(user, passwd, auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	go func() {
		s.Start()
	}()

	fmt.Println("mysql-server started!")

	<- make(chan interface{})
}

func createTestDatabase() *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()

	rows := []sql.Row{
		sql.NewRow("John Doe", "jasonkay@doe.com", []string{"555-555-555"}, time.Now()),
		sql.NewRow("John Doe", "johnalt@doe.com", []string{}, time.Now()),
		sql.NewRow("Jane Doe", "jane@doe.com", []string{}, time.Now()),
		sql.NewRow("Evil Bob", "jasonkay@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()),
	}

	for _, row := range rows {
		_ = table.Insert(ctx, row)
	}

	return db
}
