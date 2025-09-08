package service

import (
	"html/template"
	"log"
	"net/http"

	m "forum/Internal/model"
)

// HandleError renders a simple error page with the given status code.
// Falls back to http.Error if the template fails to parse/execute.
func (service *Service) HandleError(w http.ResponseWriter, status int) {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	log.Printf("http error: %d %s", status, http.StatusText(status))

	tpl, err := template.ParseFiles("./web/templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(status)
	data := m.PageData{ErrorMsg: http.StatusText(status), ErrorCode: status}
	if execErr := tpl.Execute(w, data); execErr != nil {
		http.Error(w, http.StatusText(status), status)
	}
}
