package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NOTE: Example Interaction with pgxpool lib
// func run_query(dbpool *pgxpool.Pool) string {
// 	var greeting string
// 	err := dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
//
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
// 		os.Exit(1)
// 	}
// 	return greeting
// }

// NOTE: for testing `db` package, uncomment main()
// func main() {
// 	dbpool, err := GetConnection("websitedb")
// 	defer dbpool.Close()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "GetConnection failed: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	init_posts(dbpool)
// 	view(dbpool, "posts", 4)
//
// 	fmt.Printf("We are all done!")
// }

func GetConnection(dbName string) (*pgxpool.Pool, error) {
	var (
		err error
		db  *pgxpool.Pool
	)

	if db != nil {
		return db, nil
	}
	// Init connection to PostgreSQL db
	fmt.Println("`GetConnection`: %s, %s", os.Getenv("POSTGRES_URL"), dbName)
	fmt.Println("`GetConnection`: %s", fmt.Sprintf("%s/%s", os.Getenv("POSTGRES_URL"), dbName))
	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("%s/%s", os.Getenv("POSTGRES_URL"), dbName))
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")
	return dbpool, nil
}

func init_posts(dbpool *pgxpool.Pool) {
	var (
		stmt string
		err  error
	)
	stmt = `CREATE TABLE IF NOT EXISTS posts (
		post_id 		SERIAL PRIMARY KEY,
		tags 			VARCHAR(255)[] NOT NULL,
		title 			VARCHAR(255) NOT NULL,
		link 			VARCHAR(255) NOT NULL,
		date 			DATE NOT NULL,
		content 		TEXT NOT NULL
	);`

	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_posts` failed: %v\n", err)
		os.Exit(1)
	}
}

func init_users(dbpool *pgxpool.Pool) {
	var stmt string
	var err error
	stmt = `CREATE TABLE IF NOT EXISTS posts (
		user_id 		SERIAL PRIMARY KEY,
		username 		VARCHAR(255) NOT NULL,
		password 		VARCHAR(255) NOT NULL,
	);`

	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_users` failed: %v\n", err)
		os.Exit(1)
	}
}

func GetAll(dbpool *pgxpool.Pool, table string) pgx.Rows {
	stmt := fmt.Sprintf(`SELECT * FROM %s;`, table)

	rows, err := dbpool.Query(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`GetAll` failed: %v\n", err)
		os.Exit(1)
	}
	return rows
}

func GetPosts(dbpool *pgxpool.Pool) []*Post {
	query := `SELECT * FROM posts;`

	var posts []*Post
	err := pgxscan.Select(context.Background(), dbpool, &posts, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`GetPosts` failed: %v\n", err)
		os.Exit(1)
	}
	return posts
}
