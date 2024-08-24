package controller

import (
	"html/template"
	"mvc/pkg/auth"
	"net/http"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, err := r.Cookie("token")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user.Role == "admin" {
			t := template.Must(template.ParseFiles("templates/admin_dashboard.html"))
			t.Execute(w, user)
			return
		}

		t := template.Must(template.ParseFiles("templates/user_dashboard.html"))
		t.Execute(w, user)
		return
	}
}
