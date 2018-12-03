package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join("avatars", userId+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "成功")
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	// gomniauthのセットアップ
	fbKey := os.Getenv("FB_KEY")
	fbSecretKey := os.Getenv("FB_SECRET_KEY")
	log.Println("FB_KEY: ", os.Getenv("FB_KEY"))
	log.Println("FB_SECRET_KEY: ", fbSecretKey)
	gomniauth.SetSecurityKey("securitykey")
	gomniauth.WithProviders(
		facebook.New(fbKey, fbSecretKey, "http://localhost:8080/auth/callback/facebook"),
		//github.New("クライアントID", "secretkey", "http://localhost:8080/auth/callback/github"),
		//google.New("クライアントID", "secretkey", "http://localhost:8080/auth/callback/google"),
	)
	//r := newRoom(UseAuthAvatar)
	//r := newRoom(UseGravatar)
	r := newRoom(UseFileSystemAvatar)

	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	// チャットルーム開始
	go r.run()
	// webサーバ起動
	log.Println("Webサーバを起動します。 port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndeServe:", err)
	}
}
