package impl

import (
	"context"
	"database/sql"
	f "fmt"
	i "grpc-sql-demo/server/interfaces"

	_ "github.com/go-sql-driver/mysql"

	m "grpc-sql-demo/server/models"
	pb "grpc-sql-demo/server/proto"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_grpc")
	if err != nil {
		return nil
	}

	//if db == nil {
	//	return nil
	//}
	//defer db.Close()

	return db
}

type Student struct {
	DB *sql.DB
}

func NewStudentInterfaces() i.StudentRepo {
	return &Student{
		DB: connect(),
	}
}

func (s Student) Create(ctx context.Context, in *pb.Student) (*pb.StudentId, error) {
	sql := "insert into dbgrpc values(?,?,?)"
	stmt, err := s.DB.Prepare(sql)
	if err != nil {
		f.Println("Create stmt err in Create func", err)
	}
	defer stmt.Close()

}

func (s Student) Read(ctx context.Context, in *pb.StudentId) (*pb.Student, error) {
}

func (s Student) Update(ctx context.Context, in *pb.Student) (*pb.StudentId, error) {
}

func (s Student) Delete(ctx context.Context, in *pb.StudentId) (*pb.StudentId, error) {
}

