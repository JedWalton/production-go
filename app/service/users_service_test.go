package service

import (
	"fmt"
	"math/rand"
	"os"
	"production-go/data"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func generateUniqueUsername(base string) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return fmt.Sprintf("%s_%d", base, rnd.Intn(10000))
}

func TestUserServiceWithChangePassword(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env") // Adjust the path to your .env file
	assert.NoError(t, err)

	// Initialize PostgreSQL connection
	db, err := data.NewPostgreSQL()
	assert.NoError(t, err)
	defer db.DB().Close()

	// Create a new UserService instance
	userService := NewUserService(db)

	// Define test user credentials
	username := "testuser_changepassword"
	originalPassword := "OriginalPassword123"
	newPassword := "NewPassword123"
	email := "testuser_changepassword@example.com"

	// Test user registration
	err = userService.RegisterUser(username, originalPassword, email)
	assert.NoError(t, err)

	// Ensure cleanup regardless of test result
	defer func() {
		_ = data.DeleteUser(db.DB(), username)
	}()

	// Attempt to log in with the original password
	success, err := userService.LoginUser(username, originalPassword)
	assert.NoError(t, err)
	assert.True(t, success)

	// Change the user's password
	err = userService.ChangePassword(username, originalPassword, newPassword)
	assert.NoError(t, err)

	// Attempt to log in with the old password (should fail)
	success, err = userService.LoginUser(username, originalPassword)
	assert.NoError(t, err)
	assert.False(t, success)

	// Attempt to log in with the new password (should succeed)
	success, err = userService.LoginUser(username, newPassword)
	assert.NoError(t, err)
	assert.True(t, success)
}

func TestMain(m *testing.M) {
	// Setup and teardown logic if needed
	os.Exit(m.Run())
}
