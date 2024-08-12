package main

import (
	"mvc/pkg/api"
	"mvc/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	database.Init()

	if err != nil {
		panic("Error loading .env file")
	}

	api.Start()
}
