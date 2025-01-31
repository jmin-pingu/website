package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")
	return dbpool, nil
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
