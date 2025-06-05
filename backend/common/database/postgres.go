// Package database provides functions for initializing and interacting with the database connection.
package database

import (
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver required for SQLX to connect
)

var (
	db   *sqlx.DB
	once sync.Once
)

// ConnectDB initializes and stores a PostgreSQL connection using environment config.
func ConnectDB() *sqlx.DB {
	once.Do(func() {
		dbURI := os.Getenv("DATABASE_URL")
		if dbURI == "" {
			log.Fatal("DATABASE_URL not set")
		}

		conn, err := sqlx.Connect("postgres", dbURI)
		if err != nil {
			log.Fatalf("Database connection error: %v", err)
		}

		db = conn
	})

	return db
}

// GetDB returns the initialized SQLX DB instance.
func GetDB() *sqlx.DB {
	if db == nil {
		log.Fatal("database not initialized: call ConnectDB() first")
	}
	return db
}
