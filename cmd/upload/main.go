package main

import (
	"fmt"
	"internal/db"
	"log"
	"os"
)

var VALID_RELATION = []string{"books", "posts"}

func main() {
	target := os.Args[1]
	paths := os.Args[2:]

	if target == "-h" {
		fmt.Printf("USAGE: upload [VALID_RELATION] [PATH_TO_DATA]\n\tvalid relations: %v\n", VALID_RELATION)
		return
	} else if len(paths) == 0 {
		panic("Did not provide paths to upload")
	}

	found := false
	for _, v := range VALID_RELATION {
		if target == v {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("USAGE: upload [VALID_RELATION] [PATH_TO_DATA]\n\tvalid relations: %v\n", VALID_RELATION)
		panic(fmt.Sprintf("Did not provide a valid relation."))
	}

	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	defer dbpool.Close()

	// NOTE: consider if InitPosts/InitBooks API is worth it
	// db.InitPosts(dbpool, false)
	// db.InitBooks(dbpool, false)

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
