# README

## 基本架构

作为一个简单的web服务程序，这个项目属于golang的有三部分：

- `main.go`，负责启动应用，指定监听端口
- `server/server.go`，提供后端服务
- `server/entity/jsonControl`，与存储文件`store.json`交互

