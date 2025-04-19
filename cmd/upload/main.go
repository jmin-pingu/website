package main

import (
	"fmt"
	"internal/db"
	"log"
	"os"
	"slices"
)

var VALID_RELATION = []string{"books", "posts"}

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("Not enough arguments provided")
		fmt.Printf("Usage: upload [VALID_RELATION] [PATH_TO_DATA]\n\tvalid relations: %v\n", VALID_RELATION)
		os.Exit(1)
	}

	target := os.Args[1]
	if target == "-h" {
		fmt.Printf("Usage: upload [VALID_RELATION] [PATH_TO_DATA]\n\tvalid relations: %v\n", VALID_RELATION)
		return
	}

	paths := os.Args[2:]

	if !slices.Contains(VALID_RELATION, target) {
		fmt.Printf("Usage: upload [VALID_RELATION] [PATH_TO_DATA]\n\tvalid relations: %v\n", VALID_RELATION)
		panic(fmt.Sprintf("Did not provide a valid relation."))
	}

	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	defer dbpool.Close()

	// NOTE: consider if InitPosts/InitBooks API is worth it
	db.InitPosts(dbpool)
	db.InitBooks(dbpool)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %v: %v\n", os.Getenv("POSTGRES_DB"), err)
		os.Exit(1)
	}

	switch target {
	case "books":
		db.UploadToBooks(dbpool, paths)
	case "posts":
		db.UploadToPosts(dbpool, paths)
	}
	log.Println("upload success")
}
