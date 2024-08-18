package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"mvc/pkg/auth"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"net/mail"
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
		database.DB.Model(&types.User{}).Where("admin_req <> ?", "null").Find(&users)

		err := t.Execute(w, users)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}

func MakeAdminRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
		decoded, err := auth.DecodeJWT(token.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
		id := decoded.UserID
		user := types.User{}
		err = database.DB.First(&user, id).Error
		if err != nil {
			http.Error(w, "Invalid user id!", http.StatusInternalServerError)
			return
		}
		if user.AdminReq == "pending" {
			fmt.Fprintf(w, "Admin request pending.")
			return
		} else if user.AdminReq == "denied" {
			fmt.Fprintf(w, "Admin request denied by an admin!")
			return
		}

		user.AdminReq = "pending"
		database.DB.Save(&user)

		fmt.Fprintf(w, "Admin request sent successfully!")
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

		database.DB.Model(&types.User{}).Where("id IN (?)", ids).Updates(map[string]interface{}{"role": "admin", "admin_req": "approved"})

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

		database.DB.Model(&types.User{}).Where("id IN (?)", ids).Update("admin_req", "denied")

		fmt.Fprintf(w, "%d Admin request(s) denied successfully!", len(ids))
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
		cookie, _ := r.Cookie("token")
		decoded, err := auth.DecodeJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		id := decoded.UserID
		var user types.User
		database.DB.First(&user, id)

		r.ParseForm()

		user.Username = r.FormValue("username")
		user.Email = r.FormValue("email")
		user.Phone = r.FormValue("phone")
		password := r.FormValue("password")

		pass_matched, err := auth.VerifyHash(password, user.Password)
		if err != nil || !pass_matched {
			http.Error(w, "Invalid password! Please enter correct current password!", http.StatusBadRequest)
			return
		}

		new_password := r.FormValue("newpass")
		if len(new_password) > 0 && len(new_password) < 8 {
			http.Error(w, "New password should be atleast 8 characters", http.StatusBadRequest)
			return
		} else {
			hash, err := auth.CreateHash(new_password)
			if err != nil {
				http.Error(w, "Some error occured", http.StatusInternalServerError)
				return
			}
			user.Password = hash
		}

		var temp_user types.User
		database.DB.First(&temp_user, "phone = ?", user.Phone)
		if temp_user.ID != 0 && temp_user.ID != user.ID {
			http.Error(w, "Phone number already occupied", http.StatusBadRequest)
			return
		}
		database.DB.First(&temp_user, "email = ?", user.Email)
		if temp_user.ID != 0 && temp_user.ID != user.ID {
			http.Error(w, "Email already occupied", http.StatusBadRequest)
			return
		}

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

		if len(user.Username) < 3 {
			http.Error(w, "Username too short", http.StatusBadRequest)
			return
		} else if len(user.Username) > 50 {
			http.Error(w, "Username too long", http.StatusBadRequest)
			return
		}

		_, err = mail.ParseAddress(user.Email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		user.Address = r.FormValue("address")
		if len(user.Address) < 4 {
			http.Error(w, "Address too short", http.StatusBadRequest)
			return
		}

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
