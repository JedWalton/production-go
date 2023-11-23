package data

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func generateUniqueUsername(base string) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return fmt.Sprintf("%s_%d", base, rnd.Intn(10000))
}

func TestIntegrationCreateAndDeleteUser(t *testing.T) {
	// Initialize database connection
	db, err := NewPostgreSQL()
	assert.NoError(t, err)
	defer db.DB().Close()

	// Test data
	// Use a unique username and email for the test
	username := generateUniqueUsername("testuser")
	password := "testpassword123"
	email := fmt.Sprintf("%s@example.com", username)

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// Test CreateUser
	err = CreateUser(db.DB(), username, string(passwordHash), email)
	assert.NoError(t, err)

	// Ensure cleanup regardless of test result
	defer func() {
		_ = DeleteUser(db.DB(), username)
	}()

	// Verify user creation
	var exists bool
	err = db.DB().QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Retrieve and verify password hash
	retrievedHash, err := GetPasswordHash(db.DB(), username)
	assert.NoError(t, err)
	assert.NotEmpty(t, retrievedHash)
	err = bcrypt.CompareHashAndPassword([]byte(retrievedHash), []byte(password))
	assert.NoError(t, err) // Password hash should match

	// Test DeleteUser
	err = DeleteUser(db.DB(), username)
	assert.NoError(t, err)

	// Verify user deletion
	err = db.DB().QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
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
