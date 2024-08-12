package controller

import "net/http"

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/dashboard.html")
		return
	}
}
