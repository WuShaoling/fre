# FRE Engine

Function Runtime Environment Engine

## 相关概念

- Runtime: 运行时环境，基础环境
- Template: 函数模版
- Container: 基于函数模版运行的容器(函数实例)


## 启动执行环境
docker run -v $PWD:/go/src -p 8080:80 -it --privileged golang:1.14 bash