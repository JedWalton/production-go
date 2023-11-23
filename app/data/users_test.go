package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateAndDeleteUser(t *testing.T) {
	// Initialize database connection
	db, err := NewPostgreSQL()
	assert.NoError(t, err)
	defer db.db.Close()

	// Test data
	username := "testuser"
	passwordHash := "somehash" // In real case, use a properly hashed password
	email := "testuser@example.com"

	// Test CreateUser
	err = CreateUser(db.db, username, passwordHash, email)
	assert.NoError(t, err)

	// Verify user creation
	var exists bool
	err = db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test DeleteUser
	err = DeleteUser(db.db, username)
	assert.NoError(t, err)

	// Verify user deletion
	err = db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMain(m *testing.M) {
	// Load .env file from Git root
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
	}

	os.Exit(m.Run())
}
