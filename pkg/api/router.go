package api

import (
	"fmt"
	"mvc/pkg/controller"
	"mvc/pkg/middleware"
	"net/http"
)

func Start() {
	http.Handle("/", middleware.Authenticate(http.HandlerFunc(controller.Dashboard)))
	http.Handle("/login", http.HandlerFunc(controller.Login))
	http.Handle("/register", http.HandlerFunc(controller.Register))
	http.Handle("/logout", http.HandlerFunc(controller.Logout))

	http.Handle("/books", middleware.Authenticate(http.HandlerFunc(controller.ManageBooks)))
	http.Handle("/get_books", middleware.Authenticate(http.HandlerFunc(controller.GetBooks)))
	http.Handle("/add_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.AddBook)))
	http.Handle("/remove_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.RemoveBook)))
	http.Handle("/update_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.UpdateBook)))
	http.Handle("/borrow_books", middleware.Authenticate(http.HandlerFunc(controller.BorrowBooks)))

	http.Handle("/get_users", middleware.AuthenticateAdmin(http.HandlerFunc(controller.GetUsers)))
	http.Handle("/remove_user", middleware.AuthenticateAdmin(http.HandlerFunc(controller.RemoveUser)))

	// get_book
	// add_book
	// remove_book
	// update_book

	// get_users ADMIN ONLY
	// update_user ADMIN ONLY
	// remove_user ADMIN ONLY

	http.Handle("/ping", http.HandlerFunc(controller.Ping))

	fmt.Println("Starting server on port 5050")
	http.ListenAndServe(":5050", nil)
}
