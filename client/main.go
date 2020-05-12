package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "grpc-sql-demo/client/proto"
)

const (
	port = "127.0.0.1:50053"
)

func main() {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(`Can not connect to server: %v\n`, err)
	}
	defer conn.Close()

	c := pb.NewCRUDClient(conn)

	/* Create test */
	data := &pb.Student{}
	data.Id = 8
	data.Name = "tester8"
	data.Grade = 4
	studentId, err := c.Create(context.Background(), data)
	if err != nil {
		log.Fatalf("Can't Create student, err: %v\n", err)
	}
	if studentId == nil {
		log.Fatalf("Can't Create student id: %d\n", studentId.Id)
	}
	log.Printf("Success to create student, id: %d\n", studentId.Id)

	/* Read Last Student */
	student, err := c.Read(context.Background(), &pb.StudentId{Id: 6})
	if err != nil {
		log.Fatalf("can't get data, err: %v\n", err)
	}
	log.Printf("Get student, id: %d, name: %s, grade: %d\n", student.Id,student.Name, student.Grade)

	/* Update test */
	data = &pb.Student{}
	data.Id = 7
	data.Name = "modified_tester6"
	data.Grade = 107
	updatedId, err := c.Update(context.Background(), data)
	if err != nil {
		log.Fatalf("Failed to update, err: %v\n", err)
	}
	log.Printf("Update student, id: %d\n", updatedId.Id)
}

