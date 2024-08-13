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

		isAdmin, err := auth.CheckAdmin(cookie.Value)
		if isAdmin && err == nil {
			t := template.Must(template.ParseFiles("templates/admindashboard.html"))
			t.Execute(w, nil)
			return
		}
		if !isAdmin && err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		t := template.Must(template.ParseFiles("templates/userdashboard.html"))
		t.Execute(w, nil)
		return
	}
}
