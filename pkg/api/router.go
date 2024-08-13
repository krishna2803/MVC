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

	http.Handle("/get_books", middleware.Authenticate(http.HandlerFunc(controller.GetBooks)))
	http.Handle("/books", middleware.Authenticate(http.HandlerFunc(controller.ManageBooks)))
	http.Handle("/add_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.AddBook)))
	http.Handle("/remove_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.RemoveBook)))
	http.Handle("/update_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.UpdateBook)))

	http.Handle("/get_borrows", middleware.Authenticate(http.HandlerFunc(controller.GetBorrows)))
	http.Handle("/borrows", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ManageBorrows)))
	http.Handle("/borrow_books", middleware.Authenticate(http.HandlerFunc(controller.BorrowBooks)))

	http.Handle("/get_users", middleware.AuthenticateAdmin(http.HandlerFunc(controller.GetUsers)))
	http.Handle("/remove_user", middleware.AuthenticateAdmin(http.HandlerFunc(controller.RemoveUser)))
	http.Handle("/update_user", middleware.AuthenticateSelfAndAdmin(http.HandlerFunc(controller.UpdateUser)))
	http.Handle("/users", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ManageUsers)))
	http.Handle("/admin_requests", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ManageAdminRequests)))

	http.Handle("/ping", http.HandlerFunc(controller.Ping))

	fmt.Println("Server Listening on port 5050...")
	http.ListenAndServe(":5050", nil)
}
