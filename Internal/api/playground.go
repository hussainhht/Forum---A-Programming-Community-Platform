package api

import (
	"html/template"
	"net/http"
	"strings"
	"time"

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
	tmpl, err := template.ParseFiles("./web/templates/startcoding.html") // relative path
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	data := PlaygroundPageData{
		Language: "go",
		LineEnd:  "lf",
		BOM:      "nobom",
	}
	_ = tmpl.Execute(w, data)
}

func (server *Server) Post_PlaygroundPreviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/playground", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		server.Service.HandleError(w, http.StatusBadRequest)
		return
	}

	lang := strings.ToLower(strings.TrimSpace(r.FormValue("language")))
	filename := strings.TrimSpace(r.FormValue("filename"))
	lineEnd := r.FormValue("lineend")
	bom := r.FormValue("bom")
	code := r.FormValue("code")

	if len(code) > service.MaxCodeBytes {
		http.Error(w, "Code too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Apply line endings and (optionally) BOM so preview matches the downloaded file
	code = service.ApplyLineEndings(code, lineEnd)
	codeBytes := service.ApplyBOM([]byte(code), bom == "bom")

	// Render safe highlighted HTML (your service.HighlightHTML handles both dev/audit versions)
	htmlPreview, err := service.HighlightHTML(string(codeBytes), lang)
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	data := PlaygroundPageData{
		Language:        lang,
		Filename:        filename,
		LineEnd:         lineEnd,
		BOM:             bom,
		Code:            string(codeBytes),
		HighlightedHTML: htmlPreview,
	}

	tmpl, tErr := template.ParseFiles("./web/templates/startcoding.html")
	if tErr != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
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
		server.Service.HandleError(w, http.StatusBadRequest)
		return
	}

	lang := strings.ToLower(strings.TrimSpace(r.FormValue("language")))
	filename := strings.TrimSpace(r.FormValue("filename"))
	lineEnd := r.FormValue("lineend")
	bom := r.FormValue("bom")
	code := r.FormValue("code")

	if len(code) > service.MaxCodeBytes {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
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
