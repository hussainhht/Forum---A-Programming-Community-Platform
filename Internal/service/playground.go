package service

import (
	"bytes"
	"html/template"
	"strings"

	// --- DEV ONLY (remove before audit) ---
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const MaxCodeBytes = 200_000

// ApplyLineEndings normalizes to LF then converts to CRLF if requested.
func ApplyLineEndings(s, mode string) string {
	// normalize to LF first
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	if mode == "crlf" {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	return s
}

// ApplyBOM optionally prefixes UTF-8 BOM.
func ApplyBOM(b []byte, withBOM bool) []byte {
	if !withBOM {
		return b
	}
	return append([]byte{0xEF, 0xBB, 0xBF}, b...)
}

// ResolveLang returns extension and MIME type for download.
func ResolveLang(lang string) (ext, mime string) {
	switch lang {
	case "go":
		return ".go", "text/x-go; charset=utf-8"
	case "python":
		return ".py", "text/x-python; charset=utf-8"
	case "javascript":
		return ".js", "text/javascript; charset=utf-8"
	case "typescript":
		return ".ts", "text/plain; charset=utf-8"
	case "c":
		return ".c", "text/x-c; charset=utf-8"
	case "cpp":
		return ".cpp", "text/x-c++; charset=utf-8"
	case "java":
		return ".java", "text/x-java-source; charset=utf-8"
	case "rust":
		return ".rs", "text/plain; charset=utf-8"
	case "sql":
		return ".sql", "application/sql; charset=utf-8"
	case "html":
		return ".html", "text/html; charset=utf-8"
	case "css":
		return ".css", "text/css; charset=utf-8"
	case "bash":
		return ".sh", "text/x-shellscript; charset=utf-8"
	default:
		return ".txt", "text/plain; charset=utf-8"
	}
}

// -------------------------------
// AUDIT-SAFE VERSION (no deps):
// -------------------------------
// func HighlightHTML(code, lang string) (template.HTML, error) {
// 	var buf bytes.Buffer
// 	buf.WriteString(`<pre class="chroma"><code class="language-` + stdhtml.EscapeString(lang) + `">`)
// 	buf.WriteString(stdhtml.EscapeString(code)) // escape for safety
// 	buf.WriteString(`</code></pre>`)
// 	return template.HTML(buf.String()), nil
// }

// ---------------------------------------
// DEV ONLY (remove before audit) — Chroma
// ---------------------------------------
func HighlightHTML(code, lang string) (template.HTML, error) {
	lx := lexers.Get(lang)
	if lx == nil {
		lx = lexers.Analyse(code)
	}
	if lx == nil {
		lx = lexers.Fallback
	}

	it, err := lx.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	// WithClasses(false) → inline styles (no need for extra CSS file).
	// WithClasses(true) → requires you to serve Chroma's CSS for the chosen style.
	fmtr := chromahtml.New(
		chromahtml.WithClasses(false),
		chromahtml.TabWidth(4),
	)

	st := styles.Get("monokai")
	if st == nil {
		st = styles.Fallback
	}

	if err := fmtr.Format(&b, st, it); err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil
}
