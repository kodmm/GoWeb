package main

import(
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
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
	t.templ.Execute(w, nil)
}

func main() {
	
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// WebServer start
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}