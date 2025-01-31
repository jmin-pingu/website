package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
