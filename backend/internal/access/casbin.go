// Package access provides role-based access control (RBAC) using Casbin and PostgreSQL.
package access

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var enforcer *casbin.Enforcer

// InitEnforcer initializes and returns a singleton Casbin enforcer.
func InitEnforcer(dsn string) *casbin.Enforcer {
	if enforcer != nil {
		return enforcer
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB for Casbin: %v", err)
	}

	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
	if err != nil {
		log.Fatalf("Failed to create Casbin adapter: %v", err)
	}

	// ✅ Use caller file path to resolve model.conf location safely
	_, callerPath, _, _ := runtime.Caller(0)
	modelPath := filepath.Join(filepath.Dir(callerPath), "model.conf")

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Fatalf("model.conf not found at expected path: %s", modelPath)
	}

	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	if err := e.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load Casbin policies: %v", err)
	}

	enforcer = e
	return enforcer
}

// GetEnforcer returns the already initialized Casbin enforcer.
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
