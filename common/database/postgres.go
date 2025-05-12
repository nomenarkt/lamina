// Package database provides functions for initializing and interacting with the database connection.
package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver required for SQLX to connect
)

// ConnectDB initializes and returns a PostgreSQL connection using environment config.
func ConnectDB() *sqlx.DB {
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	return db
}
