package db

import (
	"context"
	"encoding/json"
	"fmt"
	"internal/ds"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	InProgress = "In-Progress"
	ToRead     = "To-Read"
)

type Book struct {
	BookID        int
	Tags          []string
	Author        []string
	Title         string
	URL           string
	InProgress    bool
	Completed     bool
	Rating        pgtype.Float4
	DatePublished time.Time
	DateCompleted pgtype.Date
	DateStarted   pgtype.Date
}

func (b Book) Compare(other_b Book) int {
	return -1
}

type BooksJson struct {
	Books []BookJson `json:"books"`
}

type BookJson struct {
	Tags          []string `json:"tags"`
	Author        []string `json:"author"`
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	InProgress    bool     `json:"in_progress"`
	Completed     bool     `json:"completed"`
	Rating        float32  `json:"rating"`
	DatePublished string   `json:"date_published"`
	DateCompleted string   `json:"date_completed"`
	DateStarted   string   `json:"date_started"`
}

// TODO: multithread
func UploadToBooks(dbpool *pgxpool.Pool, paths []string) {
	log.Printf("--------- Uploading Paths to `books` table ---------")
	log.Printf("paths: %v", paths)

	if len(paths) != 1 {
		fmt.Fprintf(os.Stderr, "`upload.go` failed: there should be only one file for books table\n")
		os.Exit(1)
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("path '%v' exists", path)
		} else if match, _ := regexp.MatchString(".json$", path); !match {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read arguments : file should be .json\n")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read sql template: %v\n", err)
			os.Exit(1)
		}

		books := ParseBooksJson(path)
		for _, book := range books {
			tags, author, title, url, in_progress, completed, rating, date_published, date_completed, date_started := ParseBook(book)

			// cmd string, tags []string, author []string, title string, url string, in_progress bool, completed bool, rating int, date_published time.Time, date_completed time.Time)
			if Exists(dbpool, "books", "title", title) {
				UploadBook(dbpool, "update", tags, author, title, url, in_progress, completed, rating, date_published, date_completed, date_started)
			} else {
				UploadBook(dbpool, "insert", tags, author, title, url, in_progress, completed, rating, date_published, date_completed, date_started)
			}
		}
	}
}

func ParseBook(book BookJson) ([]string, []string, string, string, bool, bool, float32, time.Time, time.Time, time.Time) {
	var (
		date_published time.Time
		date_completed time.Time
		date_started   time.Time
		err            error
	)

	date_published, err = time.Parse("2006-01-02", book.DatePublished)
	if err != nil {
		log.Fatalf("failed to parse date_published from json: %v", err)
	}

	if book.InProgress {
		date_started, err = time.Parse("2006-01-02", book.DateStarted)
		if err != nil {
			log.Fatalf("failed to parse date_started from json: %v", err)
		}
	}

	if book.Completed {
		date_completed, err = time.Parse("2006-01-02", book.DateCompleted)
		if err != nil {
			log.Fatalf("failed to parse date_completed from json: %v", err)
		}
	}

	tags_set := make(ds.Set[string], 0)
	parsed_tags := book.Tags
	for _, tag := range parsed_tags {
		tags_set.Add(tag)
	}
	tags := make([]string, 0)
	for k, _ := range tags_set {
		tags = append(tags, k)
	}

	authors_set := make(ds.Set[string], 0)
	parsed_authors := book.Author
	for _, author := range parsed_authors {
		authors_set.Add(author)
	}
	authors := make([]string, 0)
	for k, _ := range authors_set {
		authors = append(authors, k)
	}
	return tags, authors, book.Title, book.URL, book.InProgress, book.Completed, book.Rating, date_published, date_completed, date_started
}

func ParseBooksJson(path string) []BookJson {
	jsonf, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to read json file: %v", err)
	}
	defer jsonf.Close()
	bytes, _ := ioutil.ReadAll(jsonf)

	var books BooksJson
	json.Unmarshal(bytes, &books)
	return books.Books
}

func InitBooks(dbpool *pgxpool.Pool, clean bool) {
	var (
		stmt string
		err  error
		dat  []byte
	)

	pwd, err := os.Getwd()
	dat, err = os.ReadFile(pwd + "/internal/db/books_schema.sql")

	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_posts` failed to read schema: %v\n", err)
		os.Exit(1)
	}

	stmt = string(dat)

	if clean {
		_, err = dbpool.Exec(context.Background(), `DROP TABLE books;`)
	}

	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_books` failed: %v\n", err)
		os.Exit(1)
	}
}

func GetBooks(dbpool *pgxpool.Pool) ds.StrictDict[string, Book] {
	var (
		books []*Book
		key   string
	)
	const QUERY = `SELECT * FROM books ORDER BY date_completed DESC, date_started DESC;`
	err := pgxscan.Select(context.Background(), dbpool, &books, QUERY)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`GetBooks` failed: %v\n", err)
		os.Exit(1)
	}

	ordered_books, _ := ds.NewStrictDict[string, Book]([]string{"2025", "2024", InProgress, ToRead})

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
			err_msg := fmt.Sprintf("failed to append to StrictDict of ordered books: \nkey: %v \nbook: %v", key, *book)
			panic(err_msg)
		}
	}
	return ordered_books
}

func UploadBook(dbpool *pgxpool.Pool, cmd string, tags []string, author []string, title string, url string, in_progress bool, completed bool, rating float32, date_published time.Time, date_completed time.Time, date_started time.Time) {
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
		parsed_rating = "'" + fmt.Sprintf("%.1f", rating) + "'"
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
		log.Printf("updated: %v", title)
	case "insert":
		script = fmt.Sprintf(INSERT_TEMPLATE, parsed_tags, parsed_authors, title, url, in_progress, completed, parsed_rating, parsed_date_published, parsed_date_completed, parsed_date_started)
		log.Printf("inserted: %v", title)

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
