package middleware

import (
	"mvc/pkg/auth"
	"net/http"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if err != nil {
			http.Error(w, "Some error occured while reading cookie", http.StatusInternalServerError)
			return
		}

		_, err = auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "token",
				Value:  "",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func AuthenticateAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if err != nil {
			http.Error(w, "Some error occured while reading cookie", http.StatusInternalServerError)
			return
		}

		decoded, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "token",
				Value:  "",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		if decoded.Role != "admin" {
			http.Error(w, "You are not authorized to access this page", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
