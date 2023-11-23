package data

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	db *sql.DB
}

func (pg *PostgreSQL) DB() *sql.DB {
	return pg.db
}

func NewPostgreSQL() (*PostgreSQL, error) {
	// Load .env file from Git root
	err := godotenv.Load("../../.env") // Adjust the relative path as necessary
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	postgresqlURL := os.Getenv("POSTGRESQL_URL")
	if postgresqlURL == "" {
		return nil, fmt.Errorf("POSTGRESQL_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", postgresqlURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreSQL{db: db}, nil
}
