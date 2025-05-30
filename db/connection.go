package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func OpenConnection() (*sql.DB, error) {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	sc := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	conn, err := sql.Open("postgres", sc)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err = conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return conn, nil
}
