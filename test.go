package main

import (
	"fmt"
	"log"
	"mywebsite/db"
	"os"
)

// This script purely exists for testing api calls + functionality + correctness
func main() {
	dbpool, err := db.GetConnection("websitedb")
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	posts := db.GetPosts(dbpool)
	fmt.Println(posts)

	for _, post := range posts {
		// TODO: fix functionality
		fmt.Println(post)
	}
}
