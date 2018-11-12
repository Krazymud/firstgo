# README

## 基本架构

作为一个简单的web服务程序，这个项目属于golang的有三部分：

- `main.go`，负责启动应用，指定监听端口
- `server/server.go`，提供后端服务
- `server/entity/jsonControl`，与存储文件`store.json`交互

除了后端的部分，前端的内容如html、css、js等内容都在public文件夹下，它们并不是这个项目的重点。

接下来就详细解析一下属于go的部分。



### main.go

在main.go中，通过pflag实现了在命令行下以参数的形式指定监听端口：

```go
pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
flag.Parse()
```

如果不特别指定端口，则监听在8080端口上。

最后，调用server.go的返回对象：

```go
server := server.NewServer()
server.Run(":" + port)
```



### server.go

在server.go中，由main.go调用的入口函数为NewServer，它的作用是添加路由等中间件，并返回一个`Negroni`实例给main.go这边来调用（通过上面的Run方法）：

```go
func NewServer() *negroni.Negroni {
	n := negroni.Classic()
	mux := http.NewServeMux()
	mux.HandleFunc("/", initial)          //初始访问
	mux.HandleFunc("/register", register) //注册
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("server/public/"))))	//静态文件
	n.UseHandler(mux)
	return n
}
```

可以看到，在这里用到了negroni这个库。Negroni是一个中间件库，它非常小，但功能很强大，它定义了中间件的框架与风格，让我们可以基于它开发出我们自己的中间件，并可以集成到Negroni中。比如我们可以把自己的业务处理Handler当作Negroni的中间件。在NewServer中，定义了一系列http mux（路由），并将这些Handler加入到Negroni中。







