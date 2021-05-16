## 在Go中集成ELK的例子

<font color="#f00">**本文建立在你已经成功构建了ELK服务基础之上；**</font>

因此，在使用本分支代码之前，请先确保已经成功部署了ELK服务；

>   如果不知道如何部署，建议先阅读：
>
>   -   Github Pages：[使用Docker-Compose部署单节点ELK](https://jasonkayzk.github.io/2021/05/15/使用Docker-Compose部署单节点ELK/)
>   -   国内Gitee镜像：[使用Docker-Compose部署单节点ELK](https://jasonkay.gitee.io/2021/05/15/使用Docker-Compose部署单节点ELK/)

此外，为了简单起见，本文中的Logstash配置和部署中的配置相同：

logstash.conf

```
input {
  tcp {
    mode => "server"
    host => "0.0.0.0"
    port => 5044
    codec => json
  }
}

output {
  elasticsearch {
    hosts => ["http://elasticsearch:9200"]
    index => "%{[service]}-%{+YYYY.MM.dd}"
  }
  stdout { codec => rubydebug }
}
```

即：

-   **LogStash通过TCP连接的方式收集日志；**
-   **同时上传ES时的索引格式为`{service}-{date}`；**

<br/>

### **在Go中使用TCP连接上传日志**

#### **编写上传代码**

既然在配置中声明的Logstash是通过TCP连接上传日志的，则我们通过在Go中创建一个TCP连接，上传日志即可；

代码如下：

logstash_demo.go

```go
package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// Logstash的TCP连接
type Logstash struct {
	Hostname   string
	Port       int
	Connection *net.TCPConn
	Timeout    int
}

// 创建一个Logstash连接
func New(hostname string, port int, timeout int) *Logstash {
	l := Logstash{}
	l.Hostname = hostname
	l.Port = port
	l.Connection = nil
	l.Timeout = timeout
	return &l
}

// 设置连接超时
func (l *Logstash) setTimeouts() {
	deadline := time.Now().Add(time.Duration(l.Timeout) * time.Millisecond)
	_ = l.Connection.SetDeadline(deadline)
	_ = l.Connection.SetWriteDeadline(deadline)
	_ = l.Connection.SetReadDeadline(deadline)
}

// 创建TCP连接
func (l *Logstash) Connect() (*net.TCPConn, error) {
	var connection *net.TCPConn
	service := fmt.Sprintf("%s:%d", l.Hostname, l.Port)
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return connection, err
	}
	connection, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return connection, err
	}
	if connection != nil {
		l.Connection = connection
		_ = l.Connection.SetLinger(0) // default -1
		_ = l.Connection.SetNoDelay(true)
		_ = l.Connection.SetKeepAlive(true)
		_ = l.Connection.SetKeepAlivePeriod(time.Duration(5) * time.Second)
		l.setTimeouts()
	}
	return connection, err
}

// 写入数据
func (l *Logstash) Writeln(message string) error {
	var err = errors.New("tpc connection is nil")
	message = fmt.Sprintf("%s\n", message)
	if l.Connection != nil {
		_, err = l.Connection.Write([]byte(message))
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				_ = l.Connection.Close()
				l.Connection = nil
			} else {
				_ = l.Connection.Close()
				l.Connection = nil
				return err
			}
		} else {
			// Successful write! Let's extend the timeout.
			l.setTimeouts()
			return nil
		}
	}
	return err
}

func main() {
	l := New("192.168.24.88", 5044, 5)
	if _, err := l.Connect(); err != nil {
		panic(err)
	}

	if err := l.Writeln(`{ "foo" : "bar", "service": "test-service" }`); err != nil {
		panic(err)
	}
}
```

代码首先创建了一个`Logstash`类，代表了一个对于Logstash的TCP连接；

函数`New`即一个初始化Logstash连接的函数；

函数`Connect`用于将当前Logstash连接对象和Logstash服务器建立连接；

函数`Writeln`用于向TCP连接中写入数据，即提交一条JSON格式的日志；

最后，在main函数中，我们首先指定logstash服务器参数并创建了一个TCP连接，随后进行了连接，并提交了一条JSON格式的日志：

```json
{ 
    "foo" : "bar", 
    "service": "test-service" 
}
```

日志中指定了`service`为`test-service`，这将通过Logstash建立一个索引`test-service-2021-05-16`的索引（因为今天是2021年05月16日）；

<br/>

#### **测试**

代码编写完毕后，接下来我们进行测试；

首先启动ELK服务：

```bash
docker-compose up -d
Creating network "elk-single_default" with the default driver
Creating elk-single_elasticsearch_1 ... done
Creating elk-single_kibana_1        ... done
Creating elk-single_logstash_1      ... done
```

访问Kibana，结果如下：

![kibana.png](https://cdn.jsdelivr.net/gh/jasonkayzk/blog_static@master/images/kibana.png)

即这时整个ELK是空的，我们没有数据，也没有为数据创建索引；

现在我们执行go项目：

```bash
go run logstash_demo.go
```

执行后查看Docker中的Logstash的日志：

```bash
docker logs -f elk-single_logstash_1
...
[2021-05-16T08:18:35,338][ERROR][logstash.inputs.tcp      ] Error in Netty pipeline: java.io.IOException: Connection reset by peer
/usr/share/logstash/vendor/bundle/jruby/2.5.0/gems/awesome_print-1.7.0/lib/awesome_print/formatters/base_formatter.rb:31: warning: constant ::Fixnum is deprecated
{
          "host" => "192.168.24.1",
           "foo" => "bar",
      "@version" => "1",
          "port" => 53635,
    "@timestamp" => 2021-05-16T08:18:35.325Z,
       "service" => "test-service"
}
```

可见我们通过TCP连接提交的日志的确被Logstash解析了；

并且刷新Kibana，可以看到已经解析到了这个索引：

![kibana_2.png](https://cdn.jsdelivr.net/gh/jasonkayzk/blog_static@master/images/kibana_2.png)

我们创建`test-service-*`的索引，并选择`Time Filter`为`@timestamp`；

随后，进行查询：

![kibana_3.png](https://cdn.jsdelivr.net/gh/jasonkayzk/blog_static@master/images/kibana_3.png)

可见，我们提交的日志的确显示在了Kibana中（忽略另外一条测试日志）；

在Go中集成ELK成功！

>   除了TCP连接之外，Logstash还支持各种各样的数据`input`形式；
>
>   这里不在介绍，感兴趣的可以看Logstash的官方文档：
>
>   -   https://www.elastic.co/guide/en/logstash/7.12/input-plugins.html

<br/>

### **使用ES Client上传日志**

除了通过Logstash对日志进行收集之外，ES本身也是支持日志提交的；

比如：通过RESTful形式的API请求提交等等；

当然ES官方也提供了Go的客户端，可以通过Go直接操作ES；

>   **既然可以通过客户端直接上传日志到ES中，为什么还要使用Logstash呢？**
>
>   这是因为Logstash中提供了大量的配置参数，可以对大量日志进行提取、过滤，并且支持各种各样的数据源；
>
>   所以在使用时，一般都会使用Logstash进行日志的过滤和整理，然后再提交至ES中；

有关ES Client，这里不再赘述，感兴趣的可以看：

-   官方仓库：https://github.com/elastic/go-elasticsearch
-   相关文章：[go-elasticsearch: Elastic 官方的 Go 语言客户端](https://www.infoq.cn/article/hvzmnkuyymckrtk-ozdp)

<br/>




