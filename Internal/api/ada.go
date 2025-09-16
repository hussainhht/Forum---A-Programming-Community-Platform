package api

import (
	"html/template"
	"net/http"
)

type AdaPageData struct {
	CountPosts   int
	Posts        interface{}
	SessionValid bool
	Username     string
	LanguageName string
	CategoryName string
	Badge        string
	Tagline      string
	MoreLink     string
	SourceLink   string
	Features     []struct{ Title, Desc string }
}

func (server *Server) Get_AdaHandler(w http.ResponseWriter, r *http.Request) {
	sessionValid := false
	username := ""
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		sessionValid = true
		if user, errU := server.Service.GetUserFromSessionID(c.Value); errU == nil && user != nil {
			username = user.Username
		}
	}

	catID, err := server.Service.GetCategoryIDByName("ada")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	posts, err := server.Service.GetPostsByCategories([]string{intToStr(catID)})
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	data := AdaPageData{
		CountPosts:   len(posts),
		Posts:        posts,
		SessionValid: sessionValid,
		Username:     username,
		LanguageName: "Ada",
		CategoryName: "ada",
		Badge:        "Ada — Safety Critical Reliability",
		Tagline:      "Structured, strongly typed, and built for mission & safety-critical software.",
		MoreLink:     "https://learn.adacore.com/",
		SourceLink:   "https://github.com/AdaCore/",
		Features: []struct{ Title, Desc string }{
			{Title: "Strong Typing", Desc: "Prevents whole classes of runtime errors at compile time."},
			{Title: "Safety & Reliability", Desc: "Built for avionics, aerospace, defense and high-assurance systems."},
			{Title: "Concurrency", Desc: "First-class tasking model with protected objects and rendezvous."},
			{Title: "Contracts & SPARK", Desc: "Formal verification via SPARK subset for mathematical proofs."},
		},
	}

	tmpl, err := template.ParseFiles("web/templates/category/Ada.html")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "Ada.html", data); err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
}
