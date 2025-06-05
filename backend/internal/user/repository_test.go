//go:build integration
// +build integration

package user

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/internal/user/testutils"
	"github.com/stretchr/testify/assert"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	testutils.RunMigrations()
	testDB = sqlx.NewDb(testutils.GetDB(), "postgres")
	m.Run()
}

func TestRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	_, err := testDB.ExecContext(ctx, `DELETE FROM users`)
	assert.NoError(t, err)

	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (
			id, email, password_hash, role, status, user_type, full_name
		) VALUES (
			42, 'test@example.com', 'hashed', 'admin', 'active', 'internal', 'Jane Doe'
		)
	`)
	assert.NoError(t, err)

	user, err := repo.FindByID(ctx, 42)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(42), user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "admin", user.Role)
	assert.Equal(t, "active", user.Status)
	assert.NotNil(t, user.FullName)
	assert.Equal(t, "Jane Doe", *user.FullName)
}

func TestRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	_, err := testDB.ExecContext(ctx, `DELETE FROM users`)
	assert.NoError(t, err)

	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (email, password_hash) VALUES 
			('a@example.com', 'a'), 
			('b@example.com', 'b')
	`)
	assert.NoError(t, err)

	users, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "a@example.com", users[0].Email)
	assert.Equal(t, "b@example.com", users[1].Email)
}

func TestRepository_UpdateUserProfile(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	// Clean slate
	_, err := testDB.ExecContext(ctx, `DELETE FROM users`)
	assert.NoError(t, err)

	// Insert baseline user
	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (
			id, email, password_hash, role, status, user_type
		) VALUES (
			101, 'update@example.com', 'hash', 'viewer', 'pending', 'external'
		)
	`)
	assert.NoError(t, err)

	fullName := "Updated User"
	companyID := 1
	phone := "1234567890"
	address := "Somewhere"

	err = repo.UpdateUserProfile(ctx, 101, fullName, &companyID, &phone, &address)
	assert.NoError(t, err)

	// Verify update
	user, err := repo.FindByID(ctx, 101)
	assert.NoError(t, err)
	assert.NotNil(t, user.FullName)
	assert.Equal(t, "Updated User", *user.FullName)
	assert.Equal(t, "1234567890", *user.Phone)
	assert.Equal(t, "Somewhere", *user.Address)
	assert.Equal(t, 1, *user.EmployeeID)
}

func TestRepository_IsAdmin(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	// Clean up and prepare fresh state
	_, err := testDB.ExecContext(ctx, `DELETE FROM users`)
	assert.NoError(t, err)

	// Insert an admin user
	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, role, status, user_type)
		VALUES (100, 'admin@example.com', 'secret', 'admin', 'active', 'internal')
	`)
	assert.NoError(t, err)

	// Insert a regular user
	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, role, status, user_type)
		VALUES (101, 'user@example.com', 'secret', 'viewer', 'active', 'external')
	`)
	assert.NoError(t, err)

	// Check admin status
	isAdmin, err := repo.IsAdmin(ctx, 100)
	assert.NoError(t, err)
	assert.True(t, isAdmin)

	isAdmin, err = repo.IsAdmin(ctx, 101)
	assert.NoError(t, err)
	assert.False(t, isAdmin)
}

func TestRepository_MarkUserActive(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	// Clean up users
	_, err := testDB.ExecContext(ctx, `DELETE FROM users`)
	assert.NoError(t, err)

	// Insert a pending user
	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, role, status, user_type)
		VALUES (200, 'pending@example.com', 'secret', 'viewer', 'pending', 'external')
	`)
	assert.NoError(t, err)

	// Call method
	err = repo.MarkUserActive(ctx, 200)
	assert.NoError(t, err)

	// Verify result
	var status string
	err = testDB.GetContext(ctx, &status, `SELECT status FROM users WHERE id = $1`, 200)
	assert.NoError(t, err)
	assert.Equal(t, "active", status)
}

func TestRepository_DeleteExpiredPendingUsers(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	// Clean the table
	_, err := testDB.ExecContext(ctx, `DELETE FROM users`)
	assert.NoError(t, err)

	// Insert users: one expired pending, one recent pending, one active
	_, err = testDB.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, role, status, user_type, created_at)
		VALUES
			(301, 'expired@example.com', 'secret', 'viewer', 'pending', 'external', NOW() - interval '25 hours'),
			(302, 'recent@example.com',  'secret', 'viewer', 'pending', 'external', NOW() - interval '1 hour'),
			(303, 'active@example.com',  'secret', 'viewer', 'active',  'external', NOW())
	`)
	assert.NoError(t, err)

	// Call the method
	err = repo.DeleteExpiredPendingUsers(ctx)
	assert.NoError(t, err)

	// Verify remaining users
	var emails []string
	err = testDB.SelectContext(ctx, &emails, `SELECT email FROM users`)
	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"recent@example.com", "active@example.com"}, emails)
}
