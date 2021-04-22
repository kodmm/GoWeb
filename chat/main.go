package main

import(
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"fmt"
	"flag"
	"github.com/kodmm/GoWeb/trace"
	"os"
)

//templは1つのテンプレートを表す
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}
//ServerHttpはHTTPリクエストを処理します。
func (t *templateHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	
	flag.Parse()
	fmt.Println(*addr)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// チャットルームの開始
	go r.run()
	// Webサーバを起動
	log.Println("Webサーバを開始する。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	
}