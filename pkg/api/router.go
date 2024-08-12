package api

import (
	"fmt"
	"mvc/pkg/controller"
	"mvc/pkg/middleware"
	"net/http"
)

func Start() {
	http.Handle("/", middleware.Authenticate(http.HandlerFunc(controller.UserDashboard)))
	http.Handle("/login", http.HandlerFunc(controller.Login))
	http.Handle("/register", http.HandlerFunc(controller.Register))
	http.Handle("/logout", http.HandlerFunc(controller.Logout))

	http.Handle("/ping", http.HandlerFunc(controller.Ping))

	fmt.Println("Starting server on port 5050")
	http.ListenAndServe(":5050", nil)
}
