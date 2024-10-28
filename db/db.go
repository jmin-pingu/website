package db

import (
	"context"
	"fmt"
	"log"
	"os"

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
	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("%s/%s", os.Getenv("POSTGRES_URL"), dbName))
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")
	return dbpool, nil
}

func InitUsers(dbpool *pgxpool.Pool) {
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

func Exists(dbpool *pgxpool.Pool, table string, column string, value string) bool {
	template := `
		SELECT 1 FROM %s
		WHERE %s='%s';
	`
	query := fmt.Sprintf(template, table, column, value)
	rows, err := dbpool.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "`Exists` failed: %v\n", err)
		os.Exit(1)
	}

	for rows.Next() {
		return true
	}
	return false
}
