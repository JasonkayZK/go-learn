# 在Docker中体验Golang的泛型

## 使用

拉取镜像并运行:

```bash
# Create Container
docker run -dit \
  --name go-v.17 \
  -v /root/workspace/go-v1.17-code:/code \
  --privileged \
  golang:1.17-rc /bin/bash
  
# Attach Container
docker exec -it go-v.17 bash
```

容器中执行Go代码：

```bash
# go run -gcflags=-G=3 1_print/main.go

1 2 3 4 5
1.01 2.02 3.03 4.04 5.05
one two three four five
5 4 3 2 1
```

## 相关博文

- Github Pages：[在Docker中体验Go1-17中的泛型](https://jasonkayzk.github.io/2021/07/05/在Docker中体验Go1-17中的泛型/)
- 国内Gitee镜像：[在Docker中体验Go1-17中的泛型](https://jasonkay.gitee.io/2021/07/05/在Docker中体验Go1-17中的泛型/)
