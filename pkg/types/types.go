package types

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"->;unique;primaryKey;autoIncrement"`
	Title    string `json:"title" gorm:"not null; type:varchar(255)"`
	Author   string `json:"author" gorm:"not null; type:varchar(255)"`
	Genre    string `json:"genre" gorm:"not null; type:varchar(255)"`
	Language string `json:"language" gorm:"not null"`
	Summary  string `json:"summary" gorm:"not null"`
	Count    int    `json:"count" gorm:"not null; default:1"`
}

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"->;unique;primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Phone    string `json:"phone" gorm:"unique; not null; type:char(10)"`
	Email    string `json:"email" gorm:"unique; not null; type:varchar(255)"`
	Address  string `json:"address" gorm:"not null"`
	Role     string `json:"role" gorm:"default:user; not null"` // user, admin
}

type Borrow struct {
	gorm.Model
	ID         uint   `json:"id" gorm:"->;unique;primaryKey;autoIncrement"`
	BookID     uint   `json:"book_id" gorm:"not null"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	Status     string `json:"status" gorm:"default:pending; not null"` // pending','approved','denied','returned
	Count      int    `json:"count" gorm:"not null"`                   // count of book when it was borrowed
	BorrowedAt string `json:"borrowed_at" gorm:"type: timestamp; default: null"`
	ReturnedAt string `json:"returned_at" gorm:"not null; type: timestamp"`
	Book       Book
	User       User
}
