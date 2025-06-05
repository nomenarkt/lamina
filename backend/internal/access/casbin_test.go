package access

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// TestPolicyEnforcement verifies that Casbin correctly enforces a basic policy and role binding.
func TestPolicyEnforcement(t *testing.T) {
	// 📍 Print working directory for debug
	dir, _ := os.Getwd()
	fmt.Println("📁 Test running from directory:", dir)

	// 🔧 Load test environment variables
	err := godotenv.Load(".env.test")
	assert.NoError(t, err, "Failed to load .env.test")

	// 📄 Build DSN from test env vars
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// 🔐 Initialize enforcer with PostgreSQL adapter (same code as production InitEnforcer)
	e := InitEnforcer(dsn)

	// 🧪 Seed test policies
	_, _ = e.AddPolicy("planner", "company:1", "/api/flights", "GET")
	_, _ = e.AddGroupingPolicy("user:42", "planner", "company:1")

	// ✅ Evaluate enforcement
	allowed, err := e.Enforce("user:42", "company:1", "/api/flights", "GET")
	assert.NoError(t, err)
	assert.True(t, allowed)
}
