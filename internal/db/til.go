package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Insert a TIL entry, or update it if one with the same link already exists
const TIL_UPSERT = `
	INSERT INTO til (til_id, tags, title, link, date, display, content)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (link) DO UPDATE
	SET tags = EXCLUDED.tags, title = EXCLUDED.title, date = EXCLUDED.date,
		display = EXCLUDED.display, content = EXCLUDED.content;
`

type TIL struct {
	TILID   string `db:"til_id"`
	Tags    []string
	Title   string
	Link    string
	Date    time.Time
	Display bool
	Content string // Format: HTML
}

func InitTIL(dbpool *pgxpool.Pool) {
	// Read SQL schema for til
	pwd, err := os.Getwd()
	dat, err := os.ReadFile(pwd + "/internal/db/til_schema.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "`InitTIL` failed to read schema: %v\n", err)
		os.Exit(1)
	}

	// Execute script
	_, err = dbpool.Exec(context.Background(), string(dat))
	if err != nil {
		fmt.Fprintf(os.Stderr, "`InitTIL` failed: %v\n", err)
		os.Exit(1)
	}
}

func GetTIL(dbpool *pgxpool.Pool) ([]*TIL, error) {
	query := `SELECT * FROM til WHERE display = 'true' ORDER BY date DESC;`

	var entries []*TIL
	err := pgxscan.Select(context.Background(), dbpool, &entries, query)
	return entries, err
}

// GetTILByLink returns nil if no displayed TIL has that link
func GetTILByLink(dbpool *pgxpool.Pool, link string) (*TIL, error) {
	query := `SELECT * FROM til WHERE link = $1 AND display = 'true';`

	var entry TIL
	err := pgxscan.Get(context.Background(), dbpool, &entry, query, link)
	if pgxscan.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// UploadToTIL returns the links it uploaded
func UploadToTIL(dbpool *pgxpool.Pool, paths []string) []string {
	links := make([]string, 0)
	log.Println("--------- Uploading Paths to `til` table ---------")
	log.Printf("paths: %v", paths)
	for _, path := range paths {
		if match, _ := regexp.MatchString(".md$", path); !match {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read arguments: file should be .md\n")
			os.Exit(1)
		}

		// Read md file and extract content and metadata (same front matter as posts)
		md, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read markdown file: %v", err)
		}
		content := mdToHTML(md)
		tags, title, link, display, date := parse_md_metadata(md)

		_, err = dbpool.Exec(context.Background(), TIL_UPSERT,
			Hash(link).String(), tags, title, link, date, display, content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "`UploadToTIL` failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\tupserted: %v\n", link)
		links = append(links, link)
	}
	return links
}
