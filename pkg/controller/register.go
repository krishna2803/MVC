package controller

import (
	"fmt"
	"html/template"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"net/mail"
	"unicode"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confpass := r.FormValue("confpass")
		address := r.FormValue("address")

		if len(password) < 8 {
			http.Error(w, "Password should be atleast 8 characters", http.StatusBadRequest)
			return
		}

		if password != confpass {
			http.Error(w, "Password and confirm password should match", http.StatusBadRequest)
			return
		}

		if len(phone) != 10 {
			http.Error(w, "Phone number should be of 10 digits", http.StatusBadRequest)
			return
		}

		for _, char := range phone {
			if !unicode.IsDigit(char) {
				http.Error(w, "Phone number should contain only digits", http.StatusBadRequest)
				return
			}
		}

		if len(username) < 3 {
			http.Error(w, "Username too short", http.StatusBadRequest)
			return
		} else if len(username) > 50 {
			http.Error(w, "Username too long", http.StatusBadRequest)
			return
		}

		_, err = mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		if len(address) < 4 {
			http.Error(w, "Address too short", http.StatusBadRequest)
			return
		}

		hash, err := auth.CreateHash(password)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		var temp_user types.User
		database.DB.First(&temp_user, "phone = ?", phone)
		if temp_user.ID != 0 {
			http.Error(w, "Phone number already taken", http.StatusBadRequest)
			return
		}
		database.DB.First(&temp_user, "email = ?", email)
		if temp_user.ID != 0 {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		user := types.User{
			Username: username,
			Phone:    phone,
			Email:    email,
			Password: hash,
			Address:  address,
		}
		database.DB.Create(&user)

		token, err := auth.CreateJWT(user)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		})

		fmt.Fprintf(w, "User successfully registered!")

		return
	}

	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/register.html"))
		t.Execute(w, nil)
		return
	}
}
