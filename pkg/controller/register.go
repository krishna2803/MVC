package controller

import (
	"fmt"
	"html/template"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confpass := r.FormValue("confpass")
		address := r.FormValue("address")

		if len(password) < 8 {
			http.Error(w, "password should be atleast 8 characters", http.StatusBadRequest)
			return
		}

		if password != confpass {
			http.Error(w, "password and confirm password should be same", http.StatusBadRequest)
			return
		}

		hash, err := auth.CreateHash(password)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		user := types.User{
			Username: username,
			Phone:    phone,
			Email:    email,
			Password: hash,
			Address:  address,
		}
		fmt.Printf("%+v\n", user)

		// save user to database
		database.DB.Create(&user)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/register.html"))
		t.Execute(w, nil)
		return
	}
}
