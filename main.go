package main

import (
	"GoArticles/database"
	"GoArticles/logging"
	"GoArticles/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logging.Log.WithError(err).Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.InitDB()
	if err != nil {
		logging.Log.WithError(err).Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	e := echo.New()

	routes.InitRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}
