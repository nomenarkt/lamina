// Package access provides Casbin access control helpers for tests.
package access

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/stretchr/testify/require"
)

// SetEnforcer injects a test enforcer into the package's global singleton
func SetEnforcer(e *casbin.Enforcer) {
	enforcer = e // uses the `var enforcer` from casbin.go
}

// InitTestEnforcer sets up an in-memory Casbin enforcer with SQLite for isolated unit testing.
func InitTestEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()

	// Resolve the absolute path to model.conf relative to this file
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)

	modelPath := filepath.Join(filepath.Dir(currentFile), "model.conf")
	_, err := os.Stat(modelPath)
	require.NoError(t, err, "model.conf not found at path: %s", modelPath)

	// Use in-memory SQLite adapter
	adapter, err := gormadapter.NewAdapter("sqlite3", ":memory:", true)
	require.NoError(t, err)

	enf, err := casbin.NewEnforcer(modelPath, adapter)
	require.NoError(t, err)

	err = enf.LoadPolicy()
	require.NoError(t, err)

	SetEnforcer(enf)
	return enf
}
