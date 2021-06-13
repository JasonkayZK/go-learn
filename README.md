# Easegress Demo

A demo to show how to use [easegress](https://github.com/megaease/easegress)(An all-rounder traffic orchestration system)

一个展示流量编排系统[easegress](https://github.com/megaease/easegress)基本使用的例子；

<br/>

## 前言

Easegress是一个开源的流量编排系统（An all-rounder traffic orchestration system），居官方介绍这个系统通过Raft共识算法（实际上就是etcd）提供了分布式情况下的高可用、可以实现流量API调度、支持高并发高性能场景；

当然在我体验之后发现其实这个系统最特别的地方（与传统Nginx相比）在于Easegress是通过插件的方式直接进行热替换进行的（如果你使用过K8S，对于这种方式应该不会陌生）；

本文介绍了Easegress的基本用法；

<br/>

## **安装Easegress**

Easegress的安装可以说是相当简单了；

其实就是两个二进制文件：

-   easegress-server
-   egctl

可以直接通过下载仓库的Release安装：

-   https://github.com/megaease/easegress/releases/tag/v1.0.0

下载后解压，然后将二进制直接放入`$PATH$`环境变量即可；

>   **当然，也可以通过编译源码安装，这里不作介绍了；**

<br/>

## **启动服务器**

启动服务器非常简单，我们可以直接使用默认参数启动：

```bash
[root@localhost easegress]# easegress-server 
WARNING: Package "github.com/golang/protobuf/protoc-gen-go/generator" is deprecated.
        A future release of golang/protobuf will delete this package,
        which has long been excluded from the compatibility promise.

2021-06-13T09:47:09.205+08:00   INFO    server/main.go:61       Easegress release: v1.0.0, repo: megaease/easegress, commit: 3dd74e5b01d0852ed5dffe03de9b4d48a64667e7
2021-06-13T09:47:09.205+08:00   INFO    storage/storage.go:250  /root/workspace/easegress/running_objects.yaml not exist
2021-06-13T09:47:09.205+08:00   INFO    cluster/cluster.go:186  starting etcd cluster
2021-06-13T09:47:09.205+08:00   INFO    cluster/cluster.go:382  client connect with endpoints: [http://localhost:2380]
2021-06-13T09:47:09.205+08:00   INFO    cluster/config.go:126   etcd config: init-cluster:eg-default-name=http://localhost:2380 cluster-state:new force-new-cluster:false
2021-06-13T09:47:09.206+08:00   INFO    cluster/cluster.go:396  client is ready
2021-06-13T09:47:09.733+08:00   INFO    cluster/cluster.go:607  server is ready
2021-06-13T09:47:09.734+08:00   INFO    cluster/cluster.go:468  lease is ready
2021-06-13T09:47:09.734+08:00   INFO    cluster/cluster.go:204  cluster is ready
2021-06-13T09:47:09.763+08:00   INFO    storage/storage.go:250  /root/workspace/easegress/running_objects.yaml not exist
2021-06-13T09:47:09.764+08:00   INFO    supervisor/supervisor.go:197    create system controller StatusSyncController
2021-06-13T09:47:09.764+08:00   INFO    cluster/cluster.go:513  session is ready
2021-06-13T09:47:09.765+08:00   INFO    api/api.go:113  api server running in localhost:2381
2021-06-13T09:47:14.738+08:00   INFO    cluster/member.go:233   self ID changed from 0 to 689e371e88f78b6a
2021-06-13T09:47:14.739+08:00   INFO    cluster/member.go:154   store clusterMembers: eg-default-name(689e371e88f78b6a)=http://localhost:2380
2021-06-13T09:47:14.739+08:00   INFO    cluster/member.go:155   store knownMembers  : eg-default-name(689e371e88f78b6a)=http://localhost:2380
```

此时会监听默认的`2379、2380、2381`三个端口；

同时会在当前启动目录创建一些文件：

-   data目录；
-   log目录；
-   member目录；
-   easegress.pid文件：记录当前easegress-server的pid；

启动成功后尝试查看成员列表：

```bash
[root@localhost easegress]# egctl member list
- options:
    name: eg-default-name
    labels: {}
    cluster-name: eg-cluster-default-name
    cluster-role: writer
    cluster-request-timeout: 10s
    cluster-listen-client-urls:
    - http://localhost:2379
    cluster-listen-peer-urls:
    - http://localhost:2380
    cluster-advertise-client-urls:
    - http://localhost:2379
    cluster-initial-advertise-peer-urls:
    - http://localhost:2380
    cluster-join-urls: []
    api-addr: localhost:2381
    debug: false
    home-dir: ./
    data-dir: data
    wal-dir: ""
    log-dir: log
    member-dir: member
    cpu-profile-file: ""
    memory-profile-file: ""
  lastHeartbeatTime: "2021-06-13T12:55:21+08:00"
  etcd:
    id: 689e371e88f78b6a
    startTime: "2021-06-13T09:47:09+08:00"
    state: Leader
```

这是我们所启动的一个单节点集群的信息；

<br/>

## **系统测试**

### **创建HTTPServer**

首先，我们创建一个`HTTPServer`来控制服务器的流量：

下面是官方提供的一个例子：

```bash
$ echo '
kind: HTTPServer
name: server-demo
port: 10080
keepAlive: true
https: false
rules:
  - paths:
    - pathPrefix: /pipeline
      backend: pipeline-demo' | egctl object create
```

上面的配置会创建一个名称为`server-demo`的`HTTPServer`类型的组件，该组件监听`10080`端口的`/pipeline`路由，并转发至`pipeline-demo`组件；

创建完成后，server端产生日志：

```bash
2021-06-13T09:48:30.195+08:00   INFO    supervisor/supervisor.go:273    create server-demo
```

>   **注意到：我们还没有创建名称叫做`pipeline-demo`的组件，因此，此时访问`127.0.0.1:10080/pipeline`会报503错误；**

如下：

```bash
[root@localhost easegress]# curl -v 127.0.0.1:10080/pipeline
* About to connect() to 127.0.0.1 port 10080 (#0)
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 10080 (#0)
> GET /pipeline HTTP/1.1
> User-Agent: curl/7.29.0
> Host: 127.0.0.1:10080
> Accept: */*
> 
< HTTP/1.1 503 Service Unavailable
< Date: Sun, 13 Jun 2021 05:08:14 GMT
< Content-Length: 0
< 
* Connection #0 to host 127.0.0.1 left intact
```

下面创建HTTPPipeline；

<br/>

### **创建HTTPPipeline**

HTTPPipeline是Easegress的核心部分，用于过滤、校验和转发HTTP请求！

下面我们来创建一个HTTPPipeline；

下面的例子同样由官方提供：

```bash
$ echo '
name: pipeline-demo
kind: HTTPPipeline
flow:
  - filter: proxy
filters:
  - name: proxy
    kind: Proxy
    mainPool:
      servers:
      - url: http://127.0.0.1:9095
      - url: http://127.0.0.1:9096
      - url: http://127.0.0.1:9097
      loadBalance:
        policy: roundRobin' | egctl object create
```

server端日志：

```bash
2021-06-13T13:14:54.667+08:00   INFO    supervisor/supervisor.go:273    create pipeline-demo
```

这个HTTPPipeline将使用`roundRobin`算法为后端的三个服务提供负载均衡：

  - url: http://127.0.0.1:9095
  - url: http://127.0.0.1:9096
  - url: http://127.0.0.1:9097

这时我们还是不能进行测试的，因为我们还没有后端服务；

所幸官方已经提供了这么一个后端服务：

-   https://github.com/megaease/easegress/blob/main/example/backend-service/mirror.go

下面我们来启动后端服务；

<br/>

### **启动后端服务**

为了让我们的测试更加明显，我对官方提供的服务做了一些修改；

修改后的后端服务的代码如下：

mirror.go

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("p", 9095, "server port, default: 9095")
	flag.Parse()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, "hello")
	}
	mirrorHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(10 * time.Millisecond)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("<read failed: %v>", err))
		}

		url := req.URL.Path
		if req.URL.Query().Encode() != "" {
			url += "?" + req.URL.Query().Encode()
		}

		content := fmt.Sprintf(`Your Request server from port %d
===============
Method: %s
URL   : %s
Header: %v
Body  : %s\n
`, *port, req.Method, url, req.Header, body)

		_, _ = io.WriteString(w, content)
	}

	http.HandleFunc("/", mirrorHandler)
	http.HandleFunc("/pipeline/activity/1", helloHandler)
	http.HandleFunc("/pipeline/activity/2", helloHandler)

	fmt.Printf("Server started at port: %d\n", *port)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
```

可以看到，服务默认监听了`9095`端口，但是我们可以通过

同时，各个路由的Handler：

-   `/`：mirrorHandler；
-   `/pipeline/activity/1`：helloHandler；
-   `/pipeline/activity/2`：helloHandler；

下面启动服务：

```bash
[root@localhost easegress]# go run mirror.go &
Server started at port: 9095
[root@localhost easegress]# go run mirror.go -p 9096 &
Server started at port: 9096
[root@localhost easegress]# go run mirror.go -p 9097 &
Server started at port: 9097

# 查看后台服务
[root@localhost easegress]# jobs
[1]   Running                 go run mirror.go &
[2]-  Running                 go run mirror.go -p 9096 &
[3]+  Running                 go run mirror.go -p 9097 &
```

服务已经成功启动，并监听`9095~9097`端口；

<br/>

### **服务测试**

注意到，在之前的配置中，我们创建的`HTTPServer`监听的是`10080`端口的`/pipeline`路由，并将流量转发至`HTTPPipeline`；

而`HTTPPipeline`会负载均衡至`http://127.0.0.1:9095`~`http://127.0.0.1:9097`三个服务；

所以我们通过访问`http://127.0.0.1:10080/pipeline`进行测试！

测试如下：

```bash
[root@localhost easegress]# curl http://127.0.0.1:10080/pipeline -d 'Hello, Easegress'
Your Request server from port 9095
===============
Method: POST
URL   : /pipeline
Header: map[Accept:[*/*] Accept-Encoding:[gzip] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.29.0]]
Body  : Hello, Easegress

[root@localhost easegress]# curl http://127.0.0.1:10080/pipeline -d 'Hello, Easegress'
Your Request server from port 9096
===============
Method: POST
URL   : /pipeline
Header: map[Accept:[*/*] Accept-Encoding:[gzip] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.29.0]]
Body  : Hello, Easegress

[root@localhost easegress]# curl http://127.0.0.1:10080/pipeline -d 'Hello, Easegress'
Your Request server from port 9097
===============
Method: POST
URL   : /pipeline
Header: map[Accept:[*/*] Accept-Encoding:[gzip] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.29.0]]
Body  : Hello, Easegress

[root@localhost easegress]# curl http://127.0.0.1:10080/pipeline -d 'Hello, Easegress'
Your Request server from port 9095
===============
Method: POST
URL   : /pipeline
Header: map[Accept:[*/*] Accept-Encoding:[gzip] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.29.0]]
Body  : Hello, Easegress
```

可以看到，服务的确是以`roundRobin`算法来进行负载均衡的！

<br/>

### **热替换Filter**

上面我们定义了一个负载均衡的Filter，具体内容如下：

```yaml
name: pipeline-demo
kind: HTTPPipeline
flow:
  - filter: proxy
filters:
  - name: proxy
    kind: Proxy
    mainPool:
      servers:
      - url: http://127.0.0.1:9095
      - url: http://127.0.0.1:9096
      - url: http://127.0.0.1:9097
      loadBalance:
        policy: roundRobin
```

下面，我们通过Easegress插件式的热替换，将现有的HTTPPipeline替换为其他类型的；

我们在原来的HTTPPipeline之上，增加参数校验和请求适配功能；

创建`pipeline-demo.yaml`文件：

pipeline-demo.yaml

```yaml
name: pipeline-demo
kind: HTTPPipeline
flow:
  - filter: validator
    jumpIf: { invalid: END }
  - filter: requestAdaptor
  - filter: proxy
filters:
  - name: validator
    kind: Validator
    headers:
      Content-Type:
        values:
        - application/json
  - name: requestAdaptor
    kind: RequestAdaptor
    header:
      set:
        X-Adapt-Key: goodplan
  - name: proxy
    kind: Proxy
    mainPool:
      servers:
      - url: http://127.0.0.1:9095
      - url: http://127.0.0.1:9096
      - url: http://127.0.0.1:9097
      loadBalance:
        policy: roundRobin
```

下面介绍一下上面的配置；

首先在`flow`部分，定义了整个HTTPPipeline的过滤流（过滤链），同时如果在`validator`部分出现了`invalid`，则直接退出Filter；

`flow`部分定义的各个步骤都在下面的`filters`数组中定义；

在`validator`中，我们定义了一个`Validator`类型的Filter，用于校验`header`中包括了`Content-Type:"application/json"`；

在`requestAdaptor`中，我们定义了`RequestAdaptor`类型的Filter，用于在响应的`Header`中添加`X-Adapt-Key: "goodplan"`；

`proxy`仍然是我们之前的定义；

文件编辑后，通过`egctl object update -f`应用（类似于`kubectl apply -f`）：

```bash
egctl object update -f pipeline-demo.yaml
```

更新后server端日志：

```bash
2021-06-13T14:19:44.233+08:00   INFO    supervisor/supervisor.go:276    update pipeline-demo
```

此时再进行测试：

```bash
[root@localhost easegress]# curl -v http://127.0.0.1:10080/pipeline -d 'Hello, Easegress'
* About to connect() to 127.0.0.1 port 10080 (#0)
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 10080 (#0)
> POST /pipeline HTTP/1.1
> User-Agent: curl/7.29.0
> Host: 127.0.0.1:10080
> Accept: */*
> Content-Length: 16
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 16 out of 16 bytes
< HTTP/1.1 400 Bad Request
< Date: Sun, 13 Jun 2021 06:20:25 GMT
< Content-Length: 0
< 
* Connection #0 to host 127.0.0.1 left intact
```

可以看到Server端返回400错误！

加上Header后再次测试：

```bash
[root@localhost easegress]# curl http://127.0.0.1:10080/pipeline -H 'Content-Type: application/json' -d '{"message": "Hello, Easegress"}'
Your Request server from port 9095
===============
Method: POST
URL   : /pipeline
Header: map[Accept:[*/*] Accept-Encoding:[gzip] Content-Type:[application/json] User-Agent:[curl/7.29.0] X-Adapt-Key:[goodplan]]
Body  : {"message": "Hello, Easegress"}
```

可以看到，请求正常返回；

同时，返回的响应中的Header部分加入了`X-Adapt-Key:[goodplan]]`；

也可以通过`egctl object list`查看相关组件：

```bash
[root@localhost easegress]# egctl object list
- filters:
  - headers:
      Content-Type:
        values:
        - application/json
    kind: Validator
    name: validator
  - header:
      set:
        X-Adapt-Key: goodplan
    kind: RequestAdaptor
    name: requestAdaptor
  - kind: Proxy
    mainPool:
      loadBalance:
        policy: roundRobin
      servers:
      - url: http://127.0.0.1:9095
      - url: http://127.0.0.1:9096
      - url: http://127.0.0.1:9097
    name: proxy
  flow:
  - filter: validator
    jumpIf:
      invalid: END
  - filter: requestAdaptor
  - filter: proxy
  kind: HTTPPipeline
  name: pipeline-demo
- https: false
  keepAlive: true
  kind: HTTPServer
  name: server-demo
  port: 10080
  rules:
  - paths:
    - backend: pipeline-demo
      pathPrefix: /pipeline
```

实验成功！

>   此处包括了更多类型的Filter说明：
>
>   -   https://github.com/megaease/easegress/blob/main/doc/filters.md

<br/>

## **总结**

从上面的实验可以看到，Easegress使用起来还是很方便的；

仅仅通过配置文件即可对请求、路由等进行设置，体验和K8S基本上类似；

当然Easegress所提供的功能还有很多很多，有兴趣的可以查看官方文档：

-   https://www.megaease.com/docs/easegress/

<br/>

## 相关文章

-   Github Pages：[流量编排系统Easegress初探](https://jasonkayzk.github.io/2021/06/13/流量编排系统Easegress初探/)
-   国内Gitee镜像：[流量编排系统Easegress初探](https://jasonkay.gitee.io/2021/06/13/流量编排系统Easegress初探/)

