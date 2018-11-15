package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	// gomniauthのセットアップa
	gomniauth.SetSecurityKey("securitykey")
	gomniauth.WithProviders(
		facebook.New(os.Getenv("FB_KEY"), os.Getenv("FB_SECRET_KEY"), "http://localhost:8080/auth/callback/facebook"),
		//github.New("クライアントID", "secretkey", "http://localhost:8080/auth/callback/github"),
		//google.New("クライアントID", "secretkey", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルーム開始
	go r.run()
	// webサーバ起動
	log.Println("Webサーバを起動します。 port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndeServe:", err)
	}
}
