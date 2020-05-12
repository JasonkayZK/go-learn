package student

import (
	"context"
	"database/sql"
	f "fmt"
	i "grpc-sql-demo/server/interfaces"

	_ "github.com/go-sql-driver/mysql"

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

func (s Student) Create(_ context.Context, in *pb.Student) (*pb.StudentId, error) {
	sql2 := "insert into dbgrpc values(?,?,?)"
	stmt, err := s.DB.Prepare(sql2)
	if err != nil {
		f.Println("Create stmt err in Create func", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		f.Printf("Failed to create student: %v, student id: %d", err, in.Id)
		return nil, err
	}

	affect, err := result.RowsAffected()
	if err != nil {
		f.Printf("Insert Student data error: %v", err)
	}
	if affect != 1 {
		f.Printf("Insert Student data error: student id:  %d", in.Id)
		return nil, err
	}
	f.Printf("Created student, id: %d, name: %s, grade: %d", in.Id, in.Name, in.Grade)

	payload := &pb.StudentId{}
	payload.Id = in.Id
	return payload, nil
}

func (s Student) Read(_ context.Context, in *pb.StudentId) (*pb.Student, error) {
	sql2 := "select id, name, grade from dbgrpc where id = ?"

	//stmt, err := s.DB.Prepare(sql2)
	//if err != nil {
	//	f.Printf("Fail to read, err: %v", err)
	//}
	//defer stmt.Close()

	//var students []m.Student
	payload := &pb.Student{}
	err := s.DB.QueryRow(sql2, in.Id).Scan(&payload.Id, &payload.Name, &payload.Grade)
	if err != nil {
		f.Printf("Fail to read, err: %v", err)
	}

	f.Printf("The student id: %d, name: %s, grade: %d", payload.Id, payload.Name, payload.Grade)
	//for rows.Next() {
	//var stu = m.Student{}
	//err := rows.Scan(&stu.Id, &stu.Name, &stu.Grade)
	//students = append(students, stu)
	//err := rows.Scan(&payload.Id, &payload.Name, &payload.Grade)
	//}

	return payload, nil
}

func (s Student) Update(_ context.Context, in *pb.Student) (*pb.StudentId, error) {
	sql2 := "update dbgrpc set name = ?, grade = ? where id = ?"

	result, err := s.DB.Exec(sql2, in.Name, in.Grade, in.Id)
	if err != nil {
		f.Printf("Fail to Update, err: %v", err)
	}

	updateId, err := result.LastInsertId()
	if err != nil {
		f.Printf("Fail to get update id, err: %v", err)
	}
	f.Printf("Id: %d's student has been update!", updateId)

	payload := &pb.StudentId{}
	payload.Id = updateId

	return payload, nil
}

func (s Student) Delete(_ context.Context, in *pb.StudentId) (*pb.StudentId, error) {
	sql2 := "delete from dbgrpc where id = ?"

	result, err := s.DB.Exec(sql2, in.Id)
	if err != nil {
		f.Printf("Fail to Delete, err: %v", err)
	}

	ids, err := result.LastInsertId()
	if err != nil {
		f.Printf("Fail to get delete id, err: %v", err)
	}
	f.Printf("Id: %d's student has been deleted!", ids)

	payload := &pb.StudentId{}
	payload.Id = ids

	return payload, nil
}
