package server

import (
	"html/template"
	"net/http"

	"github.com/Krazymud/goproject/cloudgo/server/entity"
	"github.com/urfave/negroni"
)

var prefix = "server/public/"
var urlPrefix = "/register"

func initial(w http.ResponseWriter, r *http.Request) {
	/*r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}*/
	http.Redirect(w, r, "/register", http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := r.RequestURI[len(urlPrefix):]
		if query == "" {
			t, err := template.ParseFiles(prefix + "html/register.html")
			if err != nil {
				panic(err)
			}
			err = t.Execute(w, nil)
			if err != nil {
				panic(err)
			}
		} else {
			user := entity.GetUser(query[10:])
			t, err := template.ParseFiles(prefix + "html/user.html")
			if err != nil {
				panic(err)
			}
			err = t.Execute(w, user)
			if err != nil {
				panic(err)
			}
		}
	} else {
		r.ParseForm()
		err := entity.Register(r.Form)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(200)
		}
	}
}

func NewServer() *negroni.Negroni {
	n := negroni.Classic()
	mux := http.NewServeMux()
	mux.HandleFunc("/", initial)          //设置访问的路由
	mux.HandleFunc("/register", register) //设置访问的路由
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("server/public/"))))
	n.UseHandler(mux)
	return n
}
