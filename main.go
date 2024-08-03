package main

import (
	"GoArticles/database"
	"GoArticles/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	e := echo.New()

	routes.InitRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}
