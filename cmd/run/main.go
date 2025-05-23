package main

import (
	"fmt"
	"internal/db"
	"internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to website: %v\n", err)
		os.Exit(1)
	}

	db.InitBooks(dbpool)
	db.InitPosts(dbpool)
	dbpool.Close()

	handlers.RenderPosts()
	handlers.SetUpRoutes()

	port := "8080"
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
