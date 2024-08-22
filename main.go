package main

import (
	"GoArticles/database"
	"GoArticles/logging"
	"GoArticles/routes"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error, 2) // error handling channel

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := godotenv.Load(); err != nil {
			errChan <- err
		}
	}()

	var db *sql.DB
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		db, err = database.InitDB()
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			logging.Log.WithError(err).Fatalf("Startup error: %v", err)
		}
	}

	defer db.Close()

	e := echo.New()
	routes.InitRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}
