package main

import (
	"fmt"
	"mywebsite/db"
	"os"
	"regexp"
)

func main() {
	paths := os.Args[1:]
	// TODO: add functionality with glob and paths

	// Initialize db
	dbpool, err := db.GetConnection("websitedb")
	defer dbpool.Close() // Make sure to finish the transaction

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to websitedb: %v\n", err)
		os.Exit(1)
	}

	// Parse through paths
	fmt.Println("Input paths: %s", paths)
	for _, path := range paths {
		// NOTE: a bunch of checks to make sure the file is valid
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Path is valid:", path)
		} else if match, _ := regexp.MatchString(".md$", path); !match {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read arguments : file should be .md\n")
			os.Exit(1)
		} else if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read sql template: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read sql template: %v\n", err)
			os.Exit(1)
		}

		// Read md file

		// Based on contents of md file, upload or insert
		// if db.Exists(dbpool, "posts", "link", ...) {
		// 	// TODO: call INSERT
		// } else {
		// 	// TODO: call UPDATE
		// }

	}
}
