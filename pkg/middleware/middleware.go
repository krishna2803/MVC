package middleware

import (
	"fmt"
	"mvc/pkg/auth"
	"net/http"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			fmt.Println("No cookie found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if err != nil {
			http.Error(w, "Some error occured while reading cookie", http.StatusInternalServerError)
			return
		}

		fmt.Println(cookie.Value)

		decoded, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "token",
				Value:  "",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		fmt.Println("JWT Decoded:", decoded)

		next(w, r)
	}
}
