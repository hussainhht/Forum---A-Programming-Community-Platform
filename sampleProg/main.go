package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Block struct {
	Type    string // "code" or "text"
	Content string
}

var (
	blocks []Block
	mu     sync.Mutex
)

var tpl = template.Must(template.ParseFiles("index.html"))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if r.Method == http.MethodPost {
			r.ParseForm()
			blockType := r.FormValue("type")
			content := r.FormValue("content")

			if content != "" && (blockType == "code" || blockType == "text") {
				blocks = append(blocks, Block{
					Type:    blockType,
					Content: content,
				})
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err := tpl.Execute(w, blocks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
