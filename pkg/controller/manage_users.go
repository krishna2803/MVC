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
	"unicode"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var users []types.User
		database.DB.Model(&types.User{}).Find(&users)

		json.NewEncoder(w).Encode(users)
	}
}

func ManageUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/manage_users.html"))

		var users []types.User
		database.DB.Model(&types.User{}).Find(&users)

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		for i, u := range users {
			if u.ID == user.UserID {
				users = append(users[:i], users[i+1:]...)
				break
			}
		}

		err = t.Execute(w, users)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}

func ManageAdminRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/admin_requests.html"))

		var users []types.User
		database.DB.Find(&users, "role <> ?", "null")
		database.DB.Model(&types.User{}).Where("role = ?", "pending").Find(&users)

		err := t.Execute(w, users)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}

func MakeAdminRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")
		user := types.User{}
		err := database.DB.First(&user, id).Error
		if err != nil {
			http.Error(w, "Invalid user id!", http.StatusInternalServerError)
			return
		}
		user.Role = "admin"
		user.AdminReq = "pending"
		database.DB.Save(&user)
	}
}

func ApproveAdminRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid user ids!", http.StatusInternalServerError)
			return
		}

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
				http.Error(w, "Invalid user ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		database.DB.Where("id IN (?)", ids).Updates(map[string]interface{}{"role": "admin", "admin_req": "approved"})

		fmt.Fprintf(w, "%d Admin request(s) approved successfully!", len(ids))
	}
}

func DenyAdminRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid user ids!", http.StatusInternalServerError)
			return
		}

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
				http.Error(w, "Invalid user ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		database.DB.Where("id IN (?)", ids).Update("admin_req", "denied")

		fmt.Fprintf(w, "%d Admin request(s) approved successfully!", len(ids))
	}
}

func RemoveUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")
		var str_ids []string
		err := json.Unmarshal([]byte(id), &str_ids)
		if err != nil {
			http.Error(w, "Invalid user ids!", http.StatusInternalServerError)
			return
		}

		var ids []int64
		for _, str_id := range str_ids {
			id, err := strconv.ParseInt(str_id, 10, 64)
			if err != nil || id <= 0 {
				http.Error(w, "Invalid user ids!", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		database.DB.Where("id IN (?)", ids).Delete(&types.User{})
		fmt.Fprintf(w, "%d User(s) removed successfully!", len(ids))
	}
}

func ManageUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/user_profile.html"))

		cookie, _ := r.Cookie("token")
		user, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		var user_profile types.User
		database.DB.First(&user_profile, user.UserID)

		err = t.Execute(w, user_profile)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.FormValue("id")

		var user types.User
		database.DB.First(&user, id)
		user.Username = r.FormValue("username")
		user.Email = r.FormValue("email")
		user.Phone = r.FormValue("phone")
		hash, err := auth.CreateHash(r.FormValue("password"))
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
		user.Password = hash

		if len(user.Phone) != 10 {
			http.Error(w, "Phone number should be of 10 digits", http.StatusBadRequest)
			return
		}

		for _, char := range user.Phone {
			if !unicode.IsDigit(char) {
				http.Error(w, "Phone number should contain only digits", http.StatusBadRequest)
				return
			}
		}

		if len(user.Password) < 8 {
			http.Error(w, "Password should be atleast 8 characters", http.StatusBadRequest)
			return
		}

		if user.Email == "" || user.Username == "" {
			http.Error(w, "Name and Email cannot be empty", http.StatusBadRequest)
			return
		}

		user.Address = r.FormValue("address")

		err = database.DB.Save(&user).Error
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User updated successfully!")
	}
}

func ManageHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		borrows, err := getUserBorrows(cookie)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		t := template.Must(template.ParseFiles("templates/history.html"))
		err = t.Execute(w, borrows)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}
