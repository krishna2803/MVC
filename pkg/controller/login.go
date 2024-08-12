package controller

import (
	"fmt"
	"html/template"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"unicode"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/login.html"))
		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()

		creds := r.FormValue("creds")
		password := r.FormValue("password")

		is_phone := true

		if len(creds) == 10 {
			for _, char := range creds {
				if !unicode.IsDigit(char) {
					is_phone = false
					break
				}
			}
		}

		if len(password) < 8 {
			http.Error(w, "password should be atleast 8 characters", http.StatusBadRequest)
			return
		}

		var user types.User
		if is_phone {
			database.DB.First(&user, "phone = ?", creds)
		} else {
			database.DB.First(&user, "email = ?", creds)
		}

		logged_in, err := auth.VerifyHash(password, user.Password)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		if logged_in {
			fmt.Fprintf(w, "Welcome %s!", user.Username)

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
			return
		}

		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
}
