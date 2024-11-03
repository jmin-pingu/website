package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	BookID        int
	Tags          []string
	Author        string
	Title         string
	URL           string
	InProgress    bool
	Completed     bool
	Rating        int
	DatePublished time.Time
	DateCompleted time.Time
	DateStarted   time.Time
}

func InitBooks(dbpool *pgxpool.Pool, clean bool) {
	var (
		stmt string
		err  error
		dat  []byte
	)

	pwd, err := os.Getwd()
	dat, err = os.ReadFile(pwd + "/db/books_schema.sql")

	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_posts` failed to read schema: %v\n", err)
		os.Exit(1)
	}
	stmt = string(dat)

	if clean {
		_, err = dbpool.Exec(context.Background(), `DROP TABLE books;`)
	}
	// Execute script
	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_books` failed: %v\n", err)
		os.Exit(1)
	}
}

func GetBooks(dbpool *pgxpool.Pool) []*Book {
	query := `SELECT * FROM books;`
	var books []*Book
	err := pgxscan.Select(context.Background(), dbpool, &books, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`GetBooks` failed: %v\n", err)
		os.Exit(1)
	}
	return books
}

func UploadBook(dbpool *pgxpool.Pool, cmd string, tags []string, author []string, title string, url string, in_progress bool, completed bool, rating int, date_published time.Time, date_completed time.Time, date_started time.Time) {
	var (
		script string
		err    error
	)

	UPDATE_TEMPLATE := `
		UPDATE books 
		SET tags=%v, author=%v, title='%v', url='%v', in_progress='%v', completed='%v', rating=%v, date_published='%v', date_completed=%v, date_started=%v
		WHERE title = '%v';
	`

	INSERT_TEMPLATE := `
		INSERT INTO books (tags, author, title, url, in_progress, completed, rating, date_published, date_completed, date_started)
		VALUES (
			%v,
			%v,
			'%v',
			'%v',
			'%v',
			'%v',
			%v,
			'%v',
			%v,
			%v
		);
	`

	// Parse date and tags to make sure inputs work with SQL
	parsed_date_published := strings.Split(date_published.String(), " ")[0]             // Change date to proper formatting for SQL
	parsed_date_completed := "'" + strings.Split(date_completed.String(), " ")[0] + "'" // Change date to proper formatting for SQL
	parsed_date_started := "'" + strings.Split(date_started.String(), " ")[0] + "'"     // Change date to proper formatting for SQL

	parsed_tags := "ARRAY["
	for _, v := range tags {
		parsed_tags = parsed_tags + "'" + strings.ToTitle(v) + "',"
	}
	parsed_tags = parsed_tags[:len(parsed_tags)-1] + "]"

	parsed_authors := "ARRAY["
	for _, v := range author {
		parsed_authors = parsed_authors + "'" + strings.ToTitle(v) + "',"
	}
	parsed_authors = parsed_authors[:len(parsed_authors)-1] + "]"

	// conditions to double check formatting
	var parsed_rating string
	if rating == 0 {
		parsed_rating = "NULL"
	} else {
		parsed_rating = "'" + strconv.Itoa(rating) + "'"
	}

	if !completed {
		parsed_date_completed = "NULL"
	}

	if !in_progress {
		parsed_date_started = "NULL"
	}

	switch cmd {
	case "update":
		// INSERT INTO books (tags, author, title, url, in_progress, completed, rating, date_published, date_completed)
		script = fmt.Sprintf(UPDATE_TEMPLATE, parsed_tags, parsed_authors, title, url, in_progress, completed, parsed_rating, parsed_date_published, parsed_date_completed, parsed_date_started, title)
		fmt.Printf("\tupdated book: %v\n", title)
	case "insert":
		script = fmt.Sprintf(INSERT_TEMPLATE, parsed_tags, parsed_authors, title, url, in_progress, completed, parsed_rating, parsed_date_published, parsed_date_completed, parsed_date_started)
		fmt.Printf("\tupdated book: %v\n", title)
	default:
		panic("`UploadPost`: `cmd` should either be `update` or `insert`")
	}

	fmt.Printf("\tscript: %v\n", script)

	// Execute script
	_, err = dbpool.Exec(context.Background(), script)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`UploadBook` failed: %v\n", err)
		os.Exit(1)
	}
}
