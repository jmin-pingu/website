package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"website/db"
	"website/ds"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgxpool"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// NOTE: need to consider whether I want to include the .md format in the db
// Additionally, we need to consider whether we want method to remove entries

func main() {
	VALID_DB := []string{"books", "posts"}

	target := os.Args[1]
	paths := os.Args[2:]

	// Sanity checks
	if len(paths) == 0 {
		panic("Did not provide paths to upload")
	}

	found := false
	for _, v := range VALID_DB {
		if target == v {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Sprintf("Did not provide a valid db. Choose from one of the following %s", VALID_DB))
	}

	// Initialize db
	dbpool, err := db.GetConnection("websitedb")
	defer dbpool.Close() // Make sure to finish the transaction

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to websitedb: %v\n", err)
		os.Exit(1)
	}

	switch target {
	case "books":
		UploadToBooks(dbpool, paths)
	case "posts":
		UploadToPosts(dbpool, paths)
	}
	fmt.Println("success 🎉")
}

type Books struct {
	Books []Book `json:"books"`
}

type Book struct {
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

func UploadToBooks(dbpool *pgxpool.Pool, paths []string) {
	// NOTE: check for appropriate thing to upload
	fmt.Println("--------- Uploading Paths to `books` table ---------")
	fmt.Println("Paths: %v", paths)

	if len(paths) != 1 {
		fmt.Fprintf(os.Stderr, "`upload.go` failed: there should be only one file for books table\n")
		os.Exit(1)
	}

	path := paths[0]

	if _, err := os.Stat(path); err == nil {
		fmt.Println("valid path:", path)
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
		if db.Exists(dbpool, "books", "title", title) {
			db.UploadBook(dbpool, "update", tags, author, title, url, in_progress, completed, rating, date_published, date_completed, date_started)
		} else {
			db.UploadBook(dbpool, "insert", tags, author, title, url, in_progress, completed, rating, date_published, date_completed, date_started)
		}
	}
}

// tags []string, author []string, title string, url string, in_progress bool, completed bool, rating int, date_published time.Time, date_completed time.Time)
func ParseBook(book Book) ([]string, []string, string, string, bool, bool, float32, time.Time, time.Time, time.Time) {
	var (
		date_published time.Time
		date_completed time.Time
		date_started   time.Time
		err            error
	)
	// Parse Books struct

	// Parse date
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
		// Parse date
		date_completed, err = time.Parse("2006-01-02", book.DateCompleted)
		if err != nil {
			log.Fatalf("failed to parse date_completed from json: %v", err)
		}
	}

	// Parse tags
	tags_set := make(ds.Set[string], 0)
	parsed_tags := book.Tags
	for _, tag := range parsed_tags {
		tags_set.Add(tag) // Per post tag set: ds.Set
	}
	tags := make([]string, 0)
	for k, _ := range tags_set {
		tags = append(tags, k)
	}

	// Parse tags
	authors_set := make(ds.Set[string], 0)
	parsed_authors := book.Author
	for _, author := range parsed_authors {
		authors_set.Add(author) // Per post tag set: ds.Set
	}
	authors := make([]string, 0)
	for k, _ := range authors_set {
		authors = append(authors, k)
	}
	return tags, authors, book.Title, book.URL, book.InProgress, book.Completed, book.Rating, date_published, date_completed, date_started
}

func ParseBooksJson(path string) []Book {
	// open json
	jsonf, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalf("failed to read json file: %v", err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonf.Close()

	// read json as a byte array.
	bytes, _ := ioutil.ReadAll(jsonf)

	// we initialize our Users array
	var books Books

	// we unmarshal bytes, the byte array
	json.Unmarshal(bytes, &books)
	return books.Books
}

func UploadToPosts(dbpool *pgxpool.Pool, paths []string) {
	// NOTE: check for appropriate thing to upload
	fmt.Println("--------- Uploading Paths to `posts` table ---------")
	fmt.Println("Paths: %v", paths)
	// Parse through paths
	for _, path := range paths {
		// NOTE: a bunch of checks to make sure the file is valid
		if _, err := os.Stat(path); err == nil {
			fmt.Println("valid path:", path)
		} else if match, _ := regexp.MatchString(".md$", path); !match {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read arguments : file should be .md\n")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read sql template: %v\n", err)
			os.Exit(1)
		}

		// Read md file and extract content and metadata
		md, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read markdown file: %v", err)
		}
		content := mdToHTML(md)
		tags, title, link, display, date := parse_md_metadata(md)

		// Based on contents of md file, upload or insert
		if db.Exists(dbpool, "posts", "link", link) {
			db.UploadPost(dbpool, "update", tags, title, link, date, display, content)
		} else {
			db.UploadPost(dbpool, "insert", tags, title, link, date, display, content)
		}
	}
}

func parse_md_metadata(md []byte) ([]string, string, string, bool, time.Time) {
	// Get metadata from markdown
	metadata := getMetadata(md) // Consider what additional metadata I would want to consider

	// Parse date
	date, err := time.Parse("2006-01-02", metadata["date"].(string))
	if err != nil {
		log.Fatalf("failed to parse date from metadata: %v", err)
	}

	// Parse tags
	tags_set := make(ds.Set[string], 0)
	parsed_tags := strings.Split(metadata["tags"].(string), ",")
	for _, tag := range parsed_tags {
		t := strings.ToLower(strings.TrimSpace(tag))
		tags_set.Add(t) // Per post tag set: ds.Set
	}
	tags := make([]string, 0)
	for k, _ := range tags_set {
		tags = append(tags, k)
	}
	return tags, metadata["title"].(string), metadata["link"].(string), metadata["display"].(bool), date
}

// ----- Markdown-to-HTML Functions  -----
func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func mdToHTML(md []byte) string {
	var buf bytes.Buffer
	custom_parser := goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			meta.Meta,
			extension.Strikethrough,
			extension.Table,
			mathjax.MathJax,
		),
	)
	if err := custom_parser.Convert([]byte(md), &buf); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	return buf.String()
}

func getMetadata(md []byte) map[string]interface{} {
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	context := parser.NewContext()
	if err := markdown.Convert([]byte(md), &buf, parser.WithContext(context)); err != nil {
		log.Fatalf("failed to extract metadata: %v", err)
	}
	metadata := meta.Get(context)
	return metadata
}
