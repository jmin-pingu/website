package db

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Folder under data/ (also the table name) -> uploader for that table
var CONTENT_FOLDERS = map[string]func(*pgxpool.Pool, []string) []string{
	"posts": UploadToPosts,
	"til":   UploadToTIL,
}

// SyncContent makes the posts and til tables match the markdown files under
// root (e.g. data/): every file is upserted, then rows with no file are
// deleted. It only takes the whole folder, since the delete needs every link.
func SyncContent(dbpool *pgxpool.Pool, root string) {
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "`SyncContent`: %v should be the data folder\n", root)
		os.Exit(1)
	}

	for folder, upload := range CONTENT_FOLDERS {
		dir := filepath.Join(root, folder)
		// A missing folder would delete every row in the table, so skip it
		if _, err := os.Stat(dir); err != nil {
			log.Printf("skipping %v: %v", folder, err)
			continue
		}

		links := upload(dbpool, markdownFiles(dir))

		query := fmt.Sprintf(`DELETE FROM %s WHERE link <> ALL($1) RETURNING link;`, folder)
		rows, err := dbpool.Query(context.Background(), query, links)
		if err != nil {
			fmt.Fprintf(os.Stderr, "`SyncContent` failed to delete from %v: %v\n", folder, err)
			os.Exit(1)
		}
		for rows.Next() {
			var link string
			rows.Scan(&link)
			fmt.Printf("\tdeleted: %v\n", link)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "`SyncContent` failed to delete from %v: %v\n", folder, err)
			os.Exit(1)
		}
	}
}

// markdownFiles lists the non-empty .md files under dir
func markdownFiles(dir string) []string {
	files := make([]string, 0)
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		// Skip empty drafts; the uploaders fail on files without front matter
		if info, err := d.Info(); err != nil || info.Size() == 0 {
			return err
		}
		files = append(files, p)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "`SyncContent` failed to read %v: %v\n", dir, err)
		os.Exit(1)
	}
	return files
}
