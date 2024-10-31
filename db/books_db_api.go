package db

import (
	"context"
	"fmt"
	"os"
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

func UploadBook(dbpool *pgxpool.Pool, cmd string, tags []string, title string, link string, date time.Time, content string) {
	var (
		script string
		err    error
	)

	// book_id 			SERIAL PRIMARY KEY,
	// tags   				VARCHAR(255)[] NOT NULL,
	// author 				VARCHAR(255)[] NOT NULL,
	// title 				VARCHAR(255) NOT NULL,
	// url 				VARCHAR(1023) NOT NULL UNIQUE,
	// in_progress 		BOOLEAN NOT NULL,
	// completed 			BOOLEAN NOT NULL,
	// rating 				NUMERIC CHECK (rating >= 0 AND rating <= 10),
	// date_published 		DATE NOT NULL,
	// date_completed 		DATE,

	UPDATE_TEMPLATE := `
		UPDATE books 
		SET tags=%s, author=%s, title='%s', url='%s', in_progress='%s', completed='%s', rating='%s', date_published='%s', date_completed='%s'
		WHERE link = '%s';
	`

	INSERT_TEMPLATE := `
		INSERT INTO books (tags, author, title, url, in_progress, completed, rating, date_published, date_completed)
		VALUES (
			%s,
			%s,
			'%s',
			'%s',
			'%s',
			'%s',
			'%s',
			'%s',
			'%s',
		);
	`

	// Parse date and tags to make sure inputs work with SQL
	parsed_date := strings.Split(date.String(), " ")[0] // Change date to proper formatting for SQL

	parsed_tags := "ARRAY["
	for _, v := range tags {
		parsed_tags = parsed_tags + "'" + strings.ToTitle(v) + "',"
	}
	parsed_tags = parsed_tags[:len(parsed_tags)-1] + "]"

	switch cmd {
	case "update":
		script = fmt.Sprintf(UPDATE_TEMPLATE, parsed_tags, title, link, parsed_date, content, link)
		fmt.Printf("\tupdated link: %v\n", link)
	case "insert":
		script = fmt.Sprintf(INSERT_TEMPLATE, parsed_tags, title, link, parsed_date, content)
		fmt.Printf("\tupdated link: %v\n", link)
	default:
		panic("`UploadPost`: `cmd` should either be `update` or `insert`")
	}
	// Execute script
	_, err = dbpool.Exec(context.Background(), script)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`UploadBook` failed: %v\n", err)
		os.Exit(1)
	}
}
