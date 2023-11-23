package service

import (
	"os"
	"production-go/data"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUserServiceIntegration(t *testing.T) {
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
	username := "testuser"
	password := "testpassword123"
	email := "testuser@example.com"

	// Test user registration
	err = userService.RegisterUser(username, password, email)
	assert.NoError(t, err)

	// Attempt to log in with the registered user
	success, err := userService.LoginUser(username, password)
	assert.NoError(t, err)
	assert.True(t, success)

	// Clean up: Delete the test user
	_, err = db.DB().Exec("DELETE FROM users WHERE username = $1", username)
	assert.NoError(t, err)
}

func TestMain(m *testing.M) {
	// Setup and teardown logic if needed
	os.Exit(m.Run())
}
