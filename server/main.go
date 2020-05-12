package main

import (
	"google.golang.org/grpc"
	"grpc-sql-demo/server/interfaces/student"
	pb "grpc-sql-demo/server/proto"
	"log"
	"net"
)

const (
	port = ":50053"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	log.Printf("Start listen port%s\n", port)

	s := grpc.NewServer()
	pb.RegisterCRUDServer(s, student.NewStudentInterfaces())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("Server Started")
}
