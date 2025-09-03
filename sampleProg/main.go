package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Block represents a single block of content (text or code)
type Block struct {
	Type    string // "code" or "text"
	Content string
}

var (
	blocks []Block      // In-memory slice to store blocks
	mu     sync.Mutex   // Mutex to protect concurrent access
)

// Template with custom functions
var tpl = template.Must(template.New("index.html").Funcs(template.FuncMap{
	"not": func(b bool) bool { return !b }, // helper for negation
}).ParseFiles("index.html"))

func main() {
	// Serve static assets like CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Homepage: Show blocks and handle new block submission
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		// Handle form submission
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

		// Data to pass into the template
		data := struct {
			Blocks    []Block
			HasBlocks bool
		}{
			Blocks:    blocks,
			HasBlocks: len(blocks) > 0,
		}

		// Render template
		err := tpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handle deleting the last block
	http.HandleFunc("/delete-last-block", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if r.Method == http.MethodPost && len(blocks) > 0 {
			blocks = blocks[:len(blocks)-1] // remove last block
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
