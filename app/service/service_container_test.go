package service

import (
	"production-go/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceContainer(t *testing.T) {
	// Initialize database connection
	db, err := data.NewPostgreSQL()
	assert.NoError(t, err)
	defer db.DB().Close()

	// Create a new ServiceContainer
	serviceContainer := NewServiceContainer(db)

	// Test UserService is initialized and functional
	assert.NotNil(t, serviceContainer.UserService)

	// Define test user credentials
	username := "serviceContainerTestUser"
	password := "testpassword123"
	email := "serviceContainerTestUser@example.com"

	// FUNCTIONAL TEST Test user registration through UserService
	err = serviceContainer.UserService.RegisterUser(username, password, email)
	assert.NoError(t, err)

	// Cleanup: Delete the test user
	_, err = db.DB().Exec("DELETE FROM users WHERE username = $1", username)
	assert.NoError(t, err)
}
