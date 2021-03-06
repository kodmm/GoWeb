package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	// "github.com/kodmm/GoWeb/trace"
	"os"

	"github.com/joho/godotenv"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// 現在アクティブなAvatarの実装
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

//templは1つのテンプレートを表す
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServerHttpはHTTPリクエストを処理します。
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")

	flag.Parse()
	fmt.Println(*addr)

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("太陽の末裔")
	gomniauth.WithProviders(
		google.New(os.Getenv("Google_ClientID"), os.Getenv("Google_SecretKey"), "http://localhost:8080/auth/callback/google"),
		github.New(os.Getenv("Github_ClientID"), os.Getenv("Github_SecretKey"), "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout/", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1, //クッキーの期限までの秒数。ゼロまたは負の数値の場合、クッキーは直ちに期限切れになる。
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	// チャットルームの開始
	go r.run()
	// Webサーバを起動
	log.Println("Webサーバを開始する。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
