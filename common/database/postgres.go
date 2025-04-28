package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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
