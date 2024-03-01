package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var (
	db    *sql.DB
	err   error
	store = sessions.NewCookieStore([]byte("secret-key"))
	tpl   *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("../templates/*.html"))
}

func main() {
	db, err = sql.Open("mysql", "devuser:123456@tcp(127.0.0.1:3306)/dl-activity")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/", loginPage)
	http.HandleFunc("/home", authMiddleware(homePage))
	http.HandleFunc("/challenges", authMiddleware(challengesPage))
	http.HandleFunc("/page3", authMiddleware(page3Page))
	http.HandleFunc("/page4", authMiddleware(page4Page))
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

// HANDLER FUNCTIONS
func loginPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "home.html", nil)
}

func challengesPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "challenges.html", nil)
}

func page3Page(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "page3.html", nil)
}

func page4Page(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "page4.html", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var dbUsername, dbPassword string
	err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
	if err != nil || dbPassword != password {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
