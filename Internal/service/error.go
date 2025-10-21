package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// ErrorPageData holds the data for rendering the terminal-style error page
type ErrorPageData struct {
	Code      int
	Title     string
	Message   string
	Subtitle  string
	CLIHint   string
	RequestID string
	Path      string
	Method    string
	Query     string
}

// HandleError renders the error page using the error.html template.
// Falls back to http.Error if template rendering fails.
func (service *Service) HandleError(w http.ResponseWriter, r *http.Request, status int) {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	log.Printf("http error: %d %s", status, http.StatusText(status))

	// Prepare response headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	// Generate request ID
	reqID := fmt.Sprintf("ts-%d", time.Now().UnixNano())

	// Determine subtitle based on status code
	subtitle := getErrorSubtitle(status)

	// Extract query string if present
	query := ""
	if r != nil && r.URL.RawQuery != "" {
		query = r.URL.RawQuery
	}

	// Extract path
	path := "/"
	if r != nil {
		path = r.URL.Path
	}

	// Extract method
	method := "GET"
	if r != nil {
		method = r.Method
	}

	// Prepare data for the error template
	data := ErrorPageData{
		Code:      status,
		Title:     http.StatusText(status),
		Message:   http.StatusText(status),
		Subtitle:  subtitle,
		CLIHint:   "help — available: home discover login report search <query>",
		RequestID: reqID,
		Path:      path,
		Method:    method,
		Query:     query,
	}

	// Try to render the error template
	tpl, err := template.ParseFiles("./web/templates/error.html")
	if err != nil {
		log.Printf("parse error.html failed: %v (will fallback to http.Error)", err)
		http.Error(w, http.StatusText(status), status)
		return
	}

	if execErr := tpl.Execute(w, data); execErr != nil {
		log.Printf("render error.html failed: %v (will fallback to http.Error)", execErr)
		http.Error(w, http.StatusText(status), status)
		return
	}
}

// getErrorSubtitle returns a user-friendly subtitle for each error code
func getErrorSubtitle(status int) string {
	switch status {
	case 404:
		return "The page you're looking for doesn't exist. It might have been moved or deleted."
	case 403:
		return "You don't have permission to access this resource."
	case 401:
		return "You need to be logged in to access this page."
	case 500:
		return "Something went wrong on our end. We're working to fix it."
	case 400:
		return "The request was invalid. Please check your input and try again."
	case 405:
		return "The HTTP method used is not allowed for this resource."
	case 413:
		return "The request payload is too large."
	default:
		return "An unexpected error occurred. Please try again later."
	}
}
