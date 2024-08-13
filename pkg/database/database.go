package database

import (
	"fmt"
	"mvc/pkg/auth"
	"mvc/pkg/types"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func connect() (*gorm.DB, error) {
	fmt.Println("Connecting to database...")

	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", db_user, db_pass, db_name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// disable annoying not found errors
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

	// superadmin
	hashed_pass, err := auth.CreateHash(os.Getenv("SUPER_PASS"))
	if err != nil {
		return err
	}

	var user types.User
	if db.Where("phone = ?", os.Getenv("SUPER_USER")).First(&user).Error != nil {
		db.Create(&types.User{
			Username: "admin",
			Password: hashed_pass,
			Phone:    os.Getenv("SUPER_USER"),
			Email:    "admin@admin.com",
			Address:  "admin",
			Role:     "admin",
		})
	}
	return nil
}
