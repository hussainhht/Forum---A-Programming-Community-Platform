package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "log"
    "net/http"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "forum.db")
    if err != nil {
        log.Fatal("Database connection error:", err)
    }

    // Create table if it doesn't exist
    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL,
        password TEXT NOT NULL
    );`
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal("Table creation failed:", err)
    }

    http.HandleFunc("/", serveForm)
    http.HandleFunc("/register", registerHandler)

    fmt.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveForm(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("register.html"))
    tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")
    confirmPassword := r.FormValue("confirmPassword")

    if username == "" || email == "" || password == "" || confirmPassword == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    if password != confirmPassword {
        http.Error(w, "Passwords do not match", http.StatusBadRequest)
        return
    }

    if len(password) < 8 {
        http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
        return
    }

    // Insert user into the database
    _, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
    if err != nil {
        http.Error(w, "Registration failed (username may already exist)", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "🎉 Registration successful! Welcome, %s", username)
}
