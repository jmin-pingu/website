package main

import (
	"fmt"
	"mywebsite/db"
	"mywebsite/handlers"
	"os"

	//"github.com/kaleocheng/goldmark"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize db
	dbpool, err := db.GetConnection("websitedb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to websitedb: %v\n", err)
		os.Exit(1)
	}

	db.InitPosts(dbpool)

	dbpool.Close() // Make sure to finish the transaction

	// Setup routes
	handlers.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
