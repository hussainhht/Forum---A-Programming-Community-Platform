package backend

import "net/http"

func ServeFiles() {
	cssFS := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFS))
}
