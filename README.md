# **Distributed ID Generator**

Golang + MySQL 实现分布式ID生成服务

<br/>

## **特性**

* 高性能：分配ID只访问内存
* 分布式：横向扩展，理论无上限
* 高可靠：MySQL持久化，故障恢复快
* 唯一性：生成64位整形，整体递增，永不重复
* 易用性：可自定义ID起始位置，对外HTTP服务
* 可运维性：提供健康检查接口，通过负载均衡自动摘除故障节点

<br/>

## **启动项目**

### **初始化数据库**

创建数据库：

schema.sql

```mysql
CREATE DATABASE IF NOT EXISTS `id_alloc_db`;

USE `id_alloc_db`;

CREATE TABLE `segments`
(
    `app_tag`     VARCHAR(32) NOT NULL,
    `max_id`      BIGINT      NOT NULL,
    `step`        BIGINT      NOT NULL,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`app_tag`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
    COMMENT ='业务ID池';

INSERT INTO segments(`app_tag`, `max_id`, `step`)
VALUES ('test_business', 0, 100000);
```

### **修改配置**

config.json

```json
{
  "DSN": "root:123456@tcp(127.0.0.1:3306)/id_alloc_db",
  "table": "segments",
  "HttpPort": 8880,
  "HttpReadTimeout": 5000,
  "HttpWriteTimeout": 5000
}
```

修改`DSN`为你实际数据库的配置；

### **安装依赖并运行**

执行下面的命令安装依赖并启动服务：

```bash
go mod tidy && go run main.go
```

打印出Server启动信息则成功：

```bash
$ go mod tidy && go run main.go
server started at: localhoost:8880
```

<br/>

## **使用项目**

### **请求分配ID**

请求分配ID路由为，`/alloc?app_tag=<app_name>`；

下面为结果：

```bash
$ curl http://localhost:8880/alloc?app_tag=test_business
{"resp_code":0,"msg":"success","id":1}

$ curl http://localhost:8880/alloc?app_tag=test_business
{"resp_code":0,"msg":"success","id":2}
```

### **健康检查**

请求分配ID路由为，`/health?app_tag=<app_name>`；

下面为结果：

```bash
$ curl http://localhost:8880/health?app_tag=test_business
{"resp_code":0,"msg":"success","left":199996}

$ curl http://localhost:8880/health?app_tag=test_business
{"resp_code":0,"msg":"success","left":199996}
```

此时数据库中的内容：

```
mysql> select * from id_alloc_db.segments;

+---------------+--------+--------+---------------------+
| app_tag       | max_id | step   | update_time         |
+---------------+--------+--------+---------------------+
| test_business | 200000 | 100000 | 2021-06-20 13:07:23 |
+---------------+--------+--------+---------------------+
1 row in set (0.00 sec)
```

此时ID已经缓存至了200000！

<br/>

## **相关博文**

Github Pages：

-   [UUID生成算法-UUID还是snowflake](https://jasonkayzk.github.io/2020/02/09/UUID生成算法-UUID还是snowflake/)
-   [高性能分布式ID生成器实现方法总结](https://jasonkayzk.github.io/2021/06/20/高性能分布式ID生成器实现方法总结/)
-   [在Go中仅使用MySQL实现高性能分布式ID生成器](https://jasonkayzk.github.io/2021/06/20/在Go中仅使用MySQL实现高性能分布式ID生成器/)

国内Gitee镜像：

-   [UUID生成算法-UUID还是snowflake](https://jasonkay.gitee.io/2020/02/09/UUID生成算法-UUID还是snowflake/)
-   [高性能分布式ID生成器实现方法总结](https://jasonkay.gitee.io/2021/06/20/高性能分布式ID生成器实现方法总结/)
-   [在Go中仅使用MySQL实现高性能分布式ID生成器](https://jasonkay.gitee.io/2021/06/20/在Go中仅使用MySQL实现高性能分布式ID生成器/)

<br/>

## 设计参考

-   [Leaf——美团点评分布式ID生成系统](https://tech.meituan.com/MT_Leaf.html)
-   [go-id-alloc](https://github.com/owenliang/go-id-alloc)

