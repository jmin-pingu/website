package main

import (
	"fmt"
	"internal/db"
	"internal/handlers"
	"os"

	//"github.com/kaleocheng/goldmark"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize db
	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to website: %v\n", err)
		os.Exit(1)
	}

	// // Realistically, we really do not want the option to delete the DB.
	db.InitPosts(dbpool, false)
	db.InitBooks(dbpool, false)

	handlers.BOOKS = db.GetBooks(dbpool)
	dbpool.Close() // Make sure to finish the transaction

	// Setup routes
	handlers.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))

}
