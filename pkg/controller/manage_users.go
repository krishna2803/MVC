package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"unicode"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var users []types.User
		database.DB.Model(&types.User{}).Find(&users)

		json.NewEncoder(w).Encode(users)
	}
}

func ManageUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/manage_users.html"))

		var users []types.User
		database.DB.Model(&types.User{}).Find(&users)

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		for i, u := range users {
			if u.ID == user.UserID {
				users = append(users[:i], users[i+1:]...)
				break
			}
		}

		err = t.Execute(w, users)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}

func ManageAdminRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/admin_requests.html"))

		var users []types.User
		database.DB.Find(&users, "role <> ?", "null")
		database.DB.Model(&types.User{}).Where("role = ?", "pending").Find(&users)

		err := t.Execute(w, users)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
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

		var user types.User
		database.DB.First(&user, id)
		user.Username = r.FormValue("username")
		user.Email = r.FormValue("email")
		user.Phone = r.FormValue("phone")

		if len(user.Phone) != 10 {
			http.Error(w, "Phone number should be of 10 digits", http.StatusBadRequest)
			return
		}

		for _, char := range user.Phone {
			if !unicode.IsDigit(char) {
				http.Error(w, "Phone number should contain only digits", http.StatusBadRequest)
				return
			}
		}

		user.Address = r.FormValue("address")

		err := database.DB.Save(&user).Error
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}
