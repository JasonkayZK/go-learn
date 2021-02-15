package dao

import (
	"fmt"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"go-mysql-server-demo/models"
)

const (
	user      = "user"
	passwd    = "pass"
	address   = "localhost"
	port      = "13306"
	dbName    = "test"
	tableName = "pets"
)

var petDAO *PetDAO

func TestMain(m *testing.M) {
	db, err := models.InitDb(user, passwd, address, port, dbName)
	if err != nil {
		panic(err)
	}

	go initMySQL()

	petDAO = &PetDAO{DB: db}

	m.Run()
}

func initMySQL() {
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

	go s.Start()

	fmt.Println("mysql-server started!")
}

func createTestDatabase() *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "id", Type: sql.Int64, Nullable: false, Source: tableName},
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "age", Type: sql.Int64, Nullable: false, Source: tableName},
		{Name: "photo", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "ctime", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()

	rows := []sql.Row{
		sql.NewRow(1, "cat", 11, "", time.Now()),
		sql.NewRow(2, "dog", 21, "", time.Now()),
		sql.NewRow(3, "mouse", 31, "", time.Now()),
	}

	for _, row := range rows {
		_ = table.Insert(ctx, row)
	}

	return db
}

func TestPetDAO_CreatePet(t *testing.T) {
	err := petDAO.CreatePet(&models.Pet{
		Name:  "tiger",
		Age:   2,
		Photo: "haha.jpg",
		Ctime: time.Now(),
	})
	if err != nil {
		panic(err)
	}
}

func TestPetDAO_FindPetById(t *testing.T) {
	pet, err := petDAO.FindPetById(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(pet)
}

func TestPetDAO_Update(t *testing.T) {
	err := petDAO.Update(1, 99, "mouse")
	if err != nil {
		panic(err)
	}
}

func TestPetDAO_DeleteById(t *testing.T) {
	err := petDAO.DeleteById(1)
	if err != nil {
		panic(err)
	}
}
