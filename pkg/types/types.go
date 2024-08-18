package types

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title    string `json:"title" gorm:"not null; type:varchar(255)"`
	Author   string `json:"author" gorm:"not null; type:varchar(255)"`
	Genre    string `json:"genre" gorm:"not null; type:varchar(255)"`
	Language string `json:"language" gorm:"not null"`
	Summary  string `json:"summary" gorm:"not null"`
	Count    int    `json:"count" gorm:"not null; default:1"`
}

type APIBook struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Genre    string `json:"genre"`
	Language string `json:"language"`
	Summary  string `json:"summary"`
	Count    int    `json:"count"`
}

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Phone    string `json:"phone" gorm:"unique; not null; type:char(10)"`
	Email    string `json:"email" gorm:"unique; not null; type:varchar(255)"`
	Address  string `json:"address" gorm:"not null"`
	Role     string `json:"role" gorm:"default:user; not null"` // user, admin
	AdminReq string `json:"admin_req" gorm:"default:null;"`     // null, pending, approved, denied
}

type APIUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Role     string `json:"role"`
}

type Borrow struct {
	gorm.Model
	BookID     uint      `json:"book_id" gorm:"not null"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Status     string    `json:"status" gorm:"default:pending; not null"` // pending','approved','denied','returned
	Count      int       `json:"count" gorm:"not null"`                   // count of book when it was borrowed
	BorrowedAt time.Time `json:"borrowed_at" gorm:"type: timestamp; default: null"`
	ReturnedAt time.Time `json:"returned_at" gorm:"type: timestamp; default: null"`
	Book       Book
	User       User
}

type APIBorrow struct {
	ID         uint      `json:"id"`
	BookID     uint      `json:"book_id"`
	UserID     uint      `json:"user_id"`
	Status     string    `json:"status"`
	Count      int       `json:"count"`
	BorrowedAt time.Time `json:"borrowed_at"`
	ReturnedAt time.Time `json:"returned_at"`
}
