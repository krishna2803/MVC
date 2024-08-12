package database

import (
	"fmt"
	"mvc/pkg/types"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connect() (*gorm.DB, error) {
	fmt.Println("Connecting to database...")

	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", db_user, db_pass, db_name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database!")

	db.AutoMigrate(&types.Book{}, &types.User{}, &types.Borrow{})

	return db, nil
}

func Init() error {
	db, err := connect()
	if err != nil {
		return err
	}
	DB = db
	return nil
}
