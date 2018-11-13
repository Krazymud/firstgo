# README

## 基本架构

作为一个简单的web服务程序，这个项目属于golang的有三部分：

- `main.go`，负责启动应用，指定监听端口
- `server/server.go`，提供后端服务
- `server/entity/jsonControl`，与存储文件`store.json`交互

除了后端的部分，前端的内容如html、css、js等内容都在public文件夹下，后端与前端组合，实现了一个可供注册并查看用户信息的web程序（尽管前端并不是这个项目的重点）。

接下来就详细解析一下属于后端的部分。



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

server.go真正是后端的实现了，因为这次完成的只是一个简单的web程序，官方的`net/http`包已经够用（且十分方便），所以这里就没有选择框架了。在server.go中，由main.go调用的入口函数为NewServer，它的作用是添加路由等中间件，并返回一个`Negroni`实例给main.go这边来调用（通过上面的Run方法）：

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

注意到第一行的`n := negroni.Classic()`，这种方式会默认先添加3个中间件：

```go
func Classic() *Negroni {
	return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))
}
```

分别为Panic Recovery，Log和静态文件服务器。

接下来就看到我们设置的路由，首先是initial，它的匹配字串是"/"，也就是root页面，这里把初始界面设置为注册页面，所以这里直接重定向：

```go
func initial(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/register", http.StatusFound)
}
```

register页面是主页面，在这个页面上会有两种操作：重定向时以及直接访问时——GET，点击注册按钮时——POST，这两种操作需要区分，好在go html包已经提供了非常便捷的区分方式：

```go
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//处理GET请求
	} else {
		//处理POST请求
	}
}
```

其中，在处理GET请求时也有两种区别：

1. 请求url为"/register"，即重定向或直接访问页面时，此时请求register.html
2. 请求url为"/register?username=xxx"，即注册成功后，GET用户详情时，此时请求user.html

在获取html页面时，go也提供了很方便的template供渲染时改变元素的值，它在使用时需要在html文件中写成`{{.Object}}`的形式，其中的Object在渲染时可以替换，比如在user.html中有：

```html
<p>用户名：{{.Username}}</p>
<p>学号：{{.Number}}</p>
<p>电话：{{.Phone}}</p>
<p>邮箱：{{.Email}}</p>
```

其中的四个{{}}就是渲染时要替换的部分，而在处理时有user结构为：

```go
type user struct {
	Username string
	Email    string
	Phone    string
	Number   string
}
```

可见此结构与上面的渲染对象对应，那么在渲染时就能够很方便的替换了：

```go
t, err := template.ParseFiles(prefix + "html/user.html")
err = t.Execute(w, user)
```

ParseFiles函数解析了模板文件，而Execute函数则将需要替换的部分替换为user对象中的值，这样就很方便的得到了我们所需的html内容。

POST部分则要调用entity文件夹下的jsonControl.go来与json存储文件进行交互了（检测用户重复等）。在POST处理部分要视jsonControl.go的返回值来处理传给httpResponse的内容：

```go
r.ParseForm()	//拿到post内容
err := entity.Register(r.Form)
if err != nil {	//重复
	w.WriteHeader(404)
	w.Write([]byte(err.Error()))
} else {	//正确创建
	w.WriteHeader(200)
}
```

正确创建后返回状态码200，客户端就会要求转到用户详情页了，再次转交到GET请求。



### jsonControl.go

jsonControl.go的工作包括读写json存储文件、检测用户是否重复。其中，go官方的`encoding/json`包对json文件的读写提供了充分的支持。这里主要用到了Marshal与Unmarshal两个函数。其中，Marshal函数能够便捷地将保存struct的切片转为json格式，而Unmarshal函数则刚好相反，它可以将读入的json数据保存至存储struct的切片。使用方法也非常简单：

```go
bt, _ := json.Marshal(users)
...
json.Unmarshal(bt, &users)
```

其中users是保存user结构体的切片，这里需要注意的是，user结构体中的字段首字符需要大写，否则json读入与读出时找不到这些字段，如：

```go
type user struct {
	Username string
	Email    string
	Phone    string
	Number   string
}
```

这样就能够顺利的操作json文件了。我的做法是在init时先将json存储文件中的数据读出到users中，往后每次添加user时都即更改users又重写json文件。



## 运行结果

命令行指定端口运行：













