package controller

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  "",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
