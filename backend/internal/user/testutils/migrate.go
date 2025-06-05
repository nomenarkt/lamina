//go:build integration

package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once   sync.Once
	testDB *sql.DB
	dbErr  error
	dbURL  string
)

// RunMigrations ensures schema is applied once per test suite.
func RunMigrations() {
	once.Do(func() {
		dbURL = os.Getenv("TEST_DATABASE_URL")
		if dbURL == "" {
			log.Fatal("❌ TEST_DATABASE_URL not set")
		}

		fmt.Println("🚀 Running migrations...")
		cmd := exec.Command("migrate",
			"-database", dbURL,
			"-path", "../../migrations", "up",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}

		testDB, dbErr = sql.Open("postgres", dbURL)
		if dbErr != nil {
			log.Fatalf("❌ DB open failed: %v", dbErr)
		}
		if err := testDB.Ping(); err != nil {
			log.Fatalf("❌ DB ping failed: %v", err)
		}
	})
}

// GetDB returns a raw *sql.DB for direct usage (e.g., sqlx).
func GetDB() *sql.DB {
	if testDB == nil {
		log.Fatal("❌ test DB not initialized — call RunMigrations() first")
	}
	return testDB
}
