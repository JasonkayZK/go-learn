package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "protobuf_grpc_demo/proto"
)

const (
	// gRPC服务地址
	Address = "127.0.0.1:50052"
)

type helloService struct{}

// 定义服务接口的SayHello的实现方法
func (s helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	fmt.Printf("Get remote call from client, the context is: %s\n\n", ctx)

	resp.Message = "hello " + in.Name + "."
	fmt.Printf("Response msg: " + resp.Message)

	return resp, nil
}

var HelloService = helloService{}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
	}

	// 实现gRPC Server
	s := grpc.NewServer()

	// 注册helloServer为客户端提供服务
	// 内部调用了s.RegisterServer()
	pb.RegisterHelloServer(s, HelloService)

	println("Listen on: " + Address)

	_ = s.Serve(listen)
}
