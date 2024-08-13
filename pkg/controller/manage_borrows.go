package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"strconv"
	"time"
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

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
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
				http.Error(w, "Book(s) out of stock!", http.StatusInternalServerError)
				return
			}

			borrow := types.Borrow{
				BookID: uint(id),
				UserID: user.UserID,
				Count:  book.Count,
				Status: "pending",
			}

			database.DB.Create(&borrow)
		}

		fmt.Fprintf(w, "Requested %d book(s) successfully!", len(ids))
		return
	}
}

func getBorrows() []map[string]interface{} {
	var api_borrows []types.APIBorrow
	database.DB.Model(&types.Borrow{}).Find(&api_borrows)

	borrows := make([]map[string]interface{}, 0)
	for _, borrow := range api_borrows {
		var book types.APIBook
		database.DB.Model(&types.Book{}).First(&book, borrow.BookID)

		var user types.APIUser
		database.DB.Model(&types.User{}).First(&user, borrow.UserID)

		borrows = append(borrows, map[string]interface{}{
			"ID":         borrow.ID,
			"Title":      book.Title,
			"Author":     book.Author,
			"Name":       user.Username,
			"Phone":      user.Phone,
			"Email":      user.Email,
			"Address":    user.Address,
			"Count":      book.Count,
			"BorrowedAt": borrow.BorrowedAt,
			"ReturnedAt": borrow.ReturnedAt,
			"Status":     borrow.Status,
		})
	}
	return borrows
}

func GetBorrows(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		borrows := getBorrows()
		json.NewEncoder(w).Encode(borrows)
	}
}

func ManageBorrows(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		borrows := getBorrows()
		t := template.Must(template.ParseFiles("templates/manage_borrows.html"))
		t.Execute(w, borrows)
		return
	}
}

func ApproveBorrows(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid borrow ids!", http.StatusInternalServerError)
			return
		}

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
				http.Error(w, "Invalid borrow ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		for _, id := range ids {
			borrow := types.Borrow{}
			database.DB.First(&borrow, id)

			borrow.Status = "approved"
			borrow.BorrowedAt = time.Now()

			var book types.Book
			database.DB.First(&book, borrow.BookID)

			if book.Count <= 0 {
				http.Error(w, "Book(s) out of stock!", http.StatusInternalServerError)
				return
			}

			book.Count--
			database.DB.Save(&book)
			database.DB.Save(&borrow)
		}

		fmt.Fprintf(w, "%d Borrow(s) approved successfully!", len(ids))
		return
	}
}

func DenyBorrows(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid borrow ids!", http.StatusInternalServerError)
			return
		}

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
				http.Error(w, "Invalid borrow ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		for _, id := range ids {
			borrow := types.Borrow{}
			database.DB.First(&borrow, id)

			borrow.Status = "denied"
			database.DB.Save(&borrow)
		}

		fmt.Fprintf(w, "%d Borrow(s) denied successfully!", len(ids))
		return
	}
}

func ReturnBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid borrow ids!", http.StatusInternalServerError)
			return
		}

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
				http.Error(w, "Invalid borrow ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		cookie, _ := r.Cookie("token")
		user, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		for _, id := range ids {
			borrow := types.Borrow{}
			database.DB.Find(&borrow, id).Where("user_id = ?", user.UserID)

			borrow.Status = "returned"
			borrow.ReturnedAt = time.Now()

			var book types.Book
			database.DB.First(&book, borrow.BookID)

			book.Count++
			database.DB.Save(&book)
			database.DB.Save(&borrow)
		}

		fmt.Fprintf(w, "%d Books(s) returned successfully!", len(ids))
		return
	}
}
