package api

import (
	"fmt"
	"mvc/pkg/controller"
	"mvc/pkg/middleware"
	"net/http"
)

func Start() {
	http.Handle("/", middleware.Authenticate(http.HandlerFunc(controller.Dashboard)))

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// authentication
	http.Handle("/login", http.HandlerFunc(controller.Login))
	http.Handle("/register", http.HandlerFunc(controller.Register))
	http.Handle("/logout", http.HandlerFunc(controller.Logout))

	// books
	http.Handle("/get_books", middleware.Authenticate(http.HandlerFunc(controller.GetBooks)))
	http.Handle("/books", middleware.Authenticate(http.HandlerFunc(controller.ManageBooks)))
	http.Handle("/add_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.AddBook)))
	http.Handle("/remove_books", middleware.AuthenticateAdmin(http.HandlerFunc(controller.RemoveBooks)))
	http.Handle("/update_book", middleware.AuthenticateAdmin(http.HandlerFunc(controller.UpdateBook)))

	// borrowing
	http.Handle("/get_borrows", middleware.Authenticate(http.HandlerFunc(controller.GetBorrows)))
	http.Handle("/borrows", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ManageBorrows)))
	http.Handle("/borrow_books", middleware.Authenticate(http.HandlerFunc(controller.BorrowBooks)))
	http.Handle("/approve_borrows", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ApproveBorrows)))
	http.Handle("/deny_borrows", middleware.AuthenticateAdmin(http.HandlerFunc(controller.DenyBorrows)))
	http.Handle("/return_books", middleware.Authenticate(http.HandlerFunc(controller.ReturnBooks)))

	// users
	http.Handle("/get_users", middleware.AuthenticateAdmin(http.HandlerFunc(controller.GetUsers)))
	http.Handle("/remove_users", middleware.AuthenticateAdmin(http.HandlerFunc(controller.RemoveUsers)))
	http.Handle("/update_user", middleware.Authenticate(http.HandlerFunc(controller.UpdateUser)))
	http.Handle("/users", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ManageUsers)))
	http.Handle("/admin_requests", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ManageAdminRequests)))
	http.Handle("/make_admin_request", middleware.Authenticate(http.HandlerFunc(controller.MakeAdminRequest)))
	http.Handle("/approve_admin_requests", middleware.AuthenticateAdmin(http.HandlerFunc(controller.ApproveAdminRequests)))
	http.Handle("/deny_admin_requests", middleware.AuthenticateAdmin(http.HandlerFunc(controller.DenyAdminRequests)))

	// user specific
	http.Handle("/profile", middleware.Authenticate(http.HandlerFunc(controller.ManageUserProfile)))
	http.Handle("/history", middleware.Authenticate(http.HandlerFunc(controller.ManageHistory)))

	// ping
	http.Handle("/ping", http.HandlerFunc(controller.Ping))

	fmt.Println("Server Listening on port 5050...")
	http.ListenAndServe(":5050", nil)
}
