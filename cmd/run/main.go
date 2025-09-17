package main

import (
	"fmt"
	"internal/db"
	"internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to website: %v\n", err)
		os.Exit(1)
	}

	// db.InitBooks(dbpool)
	// db.InitPosts(dbpool)
	dbpool.Close()

	handlers.RenderPosts()
	handlers.SetUpRoutes()

	port := "8080"
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
