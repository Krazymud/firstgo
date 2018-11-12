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
	mux.HandleFunc("/", initial)                                                                      //初始访问
	mux.HandleFunc("/register", register)                                                             //注册
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("server/public/")))) //静态文件
	n.UseHandler(mux)
	return n
}
