## GraphQL

### 项目说明

使用Go + MySQL构建的GraphQL API；

使用到的框架：

-   [go-chi/chi](https://github.com/go-chi/chi)
-   [graphql-go/graphql](https://github.com/graphql-go/graphql)



### 启动项目

① 使用`schema.sql`创建数据库并导入数据；

② 修改main.go中的数据库配置；

③ 使用`go mod tidy`下载依赖；

④ 使用`go run main.go`启动项目；



### 测试

可以通过Postman或者curl的方式发送请求；

这是使用curl方式的结果：

```bash
$ curl --location --request POST 'http://127.0.0.1:4000/graphql' --header 'Content-Type: application/json' --data '{"query": "{users(name:\"kevin\"){id,name,age,profession,friendly}}"}'

# 返回结果
{
  "data": {
    "users": [
      {
        "age": 35,
        "friendly": true,
        "id": 1,
        "name": "kevin",
        "profession": "waiter"
      },
      {
        "age": 15,
        "friendly": true,
        "id": 5,
        "name": "kevin",
        "profession": "in school"
      }
    ]
  }
}

curl --location --request POST 'http://9.134.243.6:4000/graphql' --header 'Content-Type: application/json' --data '{"query": "{users(name:\"kevin\"){name,friendly}}"}'

# 返回结果
{
  "data": {
    "users": [
      {
        "friendly": true,
        "name": "kevin"
      },
      {
        "friendly": true,
        "name": "kevin"
      }
    ]
  }
}
```

实际上就是向url为`http://127.0.0.1:4000/graphql`的地址发送了一个Post请求，请求体为：

```json
{
    "query": "{users(name:\"kevin\"){id,name,age,profession,friendly}}"
}
```

即在query字段中是一个GraphQL请求表达式；

通过改变GraphQL表达式，我们得到了不同的响应结果！

实验成功！



### 更多说明

关于使用Go构建GraphQL API的详细说明，见：

-   Github Pages：[使用Go构建GraphQL API](https://jasonkayzk.github.io/2021/01/21/使用Go构建GraphQLAPI/)
-   Gitee Pages国内镜像：[使用Go构建GraphQL API](https://jasonkay.gitee.io/2021/01/21/使用Go构建GraphQLAPI/)

