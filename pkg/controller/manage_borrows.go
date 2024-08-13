package controller

import (
	"encoding/json"
	"fmt"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"strconv"
)

func BorrowBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid book ids!", http.StatusInternalServerError)
			return
		}

		// convert str_ids to an int array ids
		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil {
				http.Error(w, "Invalid book ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		user, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		for _, id := range ids {
			book := types.Book{}
			database.DB.First(&book, id)

			if book.Count <= 0 {
				http.Error(w, "Book out of stock!", http.StatusInternalServerError)
				return
			}

			borrow := types.Borrow{
				BookID: uint(id),
				UserID: user.UserID,
				Count:  book.Count,
			}

			book.Count--
			database.DB.Save(&book)

			database.DB.Create(&borrow)
		}

		fmt.Fprint(w, "Borrowed books successfully! Hell Yea")
		return
	}
}

func ViewBorrows(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var borrows []types.APIBorrow
		database.DB.Model(&types.Borrow{}).Find(&borrows)

		json.NewEncoder(w).Encode(borrows)
	}
}
