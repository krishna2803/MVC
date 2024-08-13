package controller

import (
	"encoding/json"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var users []types.User
		database.DB.Model(&types.User{}).Find(&users)

		json.NewEncoder(w).Encode(users)
	}
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")
		database.DB.Delete(&types.User{}, id)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")
		username := r.FormValue("username")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		address := r.FormValue("address")

		database.DB.Model(&types.User{}).Where("id = ?", id).Updates(map[string]interface{}{"username": username, "email": email, "phone": phone, "address": address})
	}
}
