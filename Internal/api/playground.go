package api

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"forum/Internal/model"
	"forum/Internal/service"
)

type PlaygroundPageData struct {
	// sticky form values
	Language string
	Filename string
	LineEnd  string // "lf" | "crlf"
	BOM      string // "nobom" | "bom"
	Code     string

	// server-generated preview
	HighlightedHTML template.HTML

	// optional error
	Error string
}

func (server *Server) Get_PlaygroundHandler(w http.ResponseWriter, r *http.Request) {
	// gather shared layout data
	sessionIDCookie, _ := r.Cookie("session_id")
	categories, _ := server.Service.GetCategories()
	isLoggedIn := false
	var user *model.User
	if sessionIDCookie != nil && server.Service.IsValidSession(sessionIDCookie.Value) {
		isLoggedIn = true
		u, err := server.Service.GetUserFromSessionID(sessionIDCookie.Value)
		if err == nil { // ignore error gracefully for playground
			user = u
		}
	}

	// base layout data
	base := model.PageData{
		IsLoggedIn: isLoggedIn,
		User:       user,
		Categories: categories,
		CSSFile:    "/web/static/css/newtyles.css",
		ExtraCSS:   []string{"/web/static/css/coding.css"},
	}

	pg := PlaygroundPageData{Language: "go", LineEnd: "lf", BOM: "nobom"}

	// compose data passed to template (embedding both structs)
	data := struct {
		model.PageData
		PlaygroundPageData
	}{base, pg}

	tmpl, err := template.ParseFiles("./web/templates/root.html", "./web/templates/startcoding.html")
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, data)
}

func (server *Server) Post_PlaygroundPreviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/playground", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		server.Service.HandleError(w, r, http.StatusBadRequest)
		return
	}

	lang := strings.ToLower(strings.TrimSpace(r.FormValue("language")))
	filename := strings.TrimSpace(r.FormValue("filename"))
	lineEnd := r.FormValue("lineend")
	bom := r.FormValue("bom")
	code := r.FormValue("code")

	if len(code) > service.MaxCodeBytes {
		server.Service.HandleError(w, r, http.StatusRequestEntityTooLarge)
		return
	}

	// Apply line endings and (optionally) BOM so preview matches the downloaded file
	code = service.ApplyLineEndings(code, lineEnd)
	codeBytes := service.ApplyBOM([]byte(code), bom == "bom")

	// Render safe highlighted HTML (your service.HighlightHTML handles both dev/audit versions)
	htmlPreview, err := service.HighlightHTML(string(codeBytes), lang)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	// layout data
	sessionIDCookie, _ := r.Cookie("session_id")
	categories, _ := server.Service.GetCategories()
	isLoggedIn := false
	var user *model.User
	if sessionIDCookie != nil && server.Service.IsValidSession(sessionIDCookie.Value) {
		isLoggedIn = true
		u, err := server.Service.GetUserFromSessionID(sessionIDCookie.Value)
		if err == nil {
			user = u
		}
	}
	base := model.PageData{
		IsLoggedIn: isLoggedIn,
		User:       user,
		Categories: categories,
		CSSFile:    "/web/static/css/newtyles.css",
		ExtraCSS:   []string{"/web/static/css/coding.css"},
	}

	pg := PlaygroundPageData{
		Language:        lang,
		Filename:        filename,
		LineEnd:         lineEnd,
		BOM:             bom,
		Code:            string(codeBytes),
		HighlightedHTML: htmlPreview,
	}

	data := struct {
		model.PageData
		PlaygroundPageData
	}{base, pg}

	tmpl, tErr := template.ParseFiles("./web/templates/root.html", "./web/templates/startcoding.html")
	if tErr != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, data)
}

func (server *Server) Post_DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/playground", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		server.Service.HandleError(w, r, http.StatusBadRequest)
		return
	}

	lang := strings.ToLower(strings.TrimSpace(r.FormValue("language")))
	filename := strings.TrimSpace(r.FormValue("filename"))
	lineEnd := r.FormValue("lineend")
	bom := r.FormValue("bom")
	code := r.FormValue("code")

	if len(code) > service.MaxCodeBytes {
		server.Service.HandleError(w, r, http.StatusRequestEntityTooLarge)
		return
	}

	ext, mime := service.ResolveLang(lang)
	if filename == "" {
		filename = "snippet_" + time.Now().Format("20060102_150405")
	}
	fullname := filename + ext

	// Apply line endings + BOM to the downloadable content
	code = service.ApplyLineEndings(code, lineEnd)
	out := service.ApplyBOM([]byte(code), bom == "bom")

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconvQuote(fullname))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

// strconvQuote ensures safe filename header formatting without importing strconv here.
func strconvQuote(s string) string {
	// minimal safe quoting for header; strip internal quotes
	return `"` + strings.ReplaceAll(s, `"`, "") + `"`
}
