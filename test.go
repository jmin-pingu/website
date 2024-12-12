package main

import (
	"context"
	"fmt"
	"mywebsite/db"
	"mywebsite/ds"
	"os"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgtype"
)

const (
	InProgress = "In-Progress"
	ToRead     = "To-Read"
)

// This script purely exists for testing api calls + functionality + correctness
func main() {
	dbpool, err := db.GetConnection("websitedb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to websitedb: %v\n", err)
		os.Exit(1)
	}

	query := `SELECT * FROM books ORDER BY date_completed;`
	var books []*db.Book
	err = pgxscan.Select(context.Background(), dbpool, &books, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`GetBooks` failed: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	ordered_books, _ := ds.NewStrictDict[string, db.Book]([]string{"2024", InProgress, ToRead})
	var key string
	for _, book := range books {
		if book.DateCompleted.Status == pgtype.Null {
			if book.InProgress {
				key = InProgress
			} else {
				key = ToRead
			}
		} else if book.DateCompleted.Status == pgtype.Present {
			key = strconv.Itoa(book.DateCompleted.Time.Year())
		} else {
			panic("Failed to parse DateCompleted for entry in books")
		}

		err = ordered_books.Append(key, *book)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println(ordered_books.Get("2024"))
}
