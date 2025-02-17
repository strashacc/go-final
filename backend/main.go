package main

import (
	"fmt"
	"os"

	"go-final/database"
	"go-final/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	database.New(os.Getenv("MONGODB_URI"))

	routes.Setup()
}


