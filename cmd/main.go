package main

import (
	"mvc/pkg/api"
	"mvc/pkg/controller"
	"mvc/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	database.Init()

	controller.AddDummyBookData()

	api.Start()
}
