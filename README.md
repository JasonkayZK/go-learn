A branch for protobuf & gRPC demo.

### 安装protoc和protoc-gen-go

#### protoc

Github的Release下载：

[https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)

下载对应平台的压缩包(win/Linux);

解压之后将bin目录下的protoc复制到系统环境变量的路径下:

-   win: `C:\Windows\System32`
-   Linux: `/bin`

等等;

然后测试:

```bash
protoc --version
libprotoc 3.12.0
```

****

#### protoc-gen-go

Github的Release下载:

[https://github.com/golang/protobuf/releases](https://github.com/golang/protobuf/releases)

下载压缩包并解压;

解压后在`protoc-gen-go`目录下编译:

```bash
go build
```

然后将生成的二进制文件`protoc-gen-go`复制到系统环境变量的路径下

全部安装完成!

<br/>

### 目录结构

```bash
tree /F                                       
.                                             
│  go.mod                                       
│  go.sum                                       
│                                               
├─hello_client                                  
│      main.go                                  
│                                               
├─hello_server                                  
│      main.go                                  
│                                               
└─proto                                         
        hello.pb.go                             
        hello.proto                             
```

<br/>

### 编写proto

本例子使用proto3语法;

hello.proto

```protobuf
// 指定proto版本
syntax = "proto3";

package proto;

// 定义请求结构
message HelloRequest{
    string name = 1;
}

// 定义响应结构
message HelloResponse{
    string message = 1;
}

// 定义Hello服务
service Hello{
    // 定义服务中的方法
    rpc SayHello (HelloRequest) returns (HelloResponse);
}
```

编写完成之后使用protoc生成`.pb.gp`文件:

```bash
protoc -I . --go_out=plugins=grpc:. ./hello.proto
```

即可生成hello.pb.go文件

<br/>

### 初始化gomod项目

使用`go mod init protobuf_grpc_demo`初始化项目;

生成go.mod文件;

然后使用`go mod tidy`初始化;

<br/>

### 服务端代码

hello_server/main.go

```go
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
```

<br/>

### 客户端代码

```go
package main

import (
   "context"
   "fmt"
   "google.golang.org/grpc"
   pb "protobuf_grpc_demo/proto"
)

const (
   Address = "127.0.0.1:50052"
)

func main() {
   // 连接gRPC服务器
   conn, err := grpc.Dial(Address, grpc.WithInsecure())
   if err != nil {
      fmt.Printf("Failed to dial to: " + Address)
   }
   if conn == nil {
      panic("Failed to get connection from the server")
   }
   defer conn.Close()

   // 初始化客户端
   c := pb.NewHelloClient(conn)

   // 调用方法
   reqBody := new(pb.HelloRequest)
   reqBody.Name = "gRPC"
   r, err := c.SayHello(context.Background(), reqBody)
   if err != nil {
      fmt.Printf("Fail to call method, err: %v", err)
   }

   fmt.Println(r.Message)
}
```

<br/>

### 运行测试

```bash
# 服务端运行：
D:\workspace\go_learn>go run hello_server/main.go
Listen on: 127.0.0.1:50052

# 客户端运行：
D:\workspace\go_learn>go run hello_client/main.go
hello gRPC.

# 并且在服务器返回
Get remote call from client, the context is: context.Background.WithCancel.WithValue(type peer.peerKey, val <not Stri
nger>).WithValue(type metadata.mdIncomingKey, val <not Stringer>).WithValue(type grpc.streamKey, val <not Stringer>)

Response msg: hello gRPC.
```

