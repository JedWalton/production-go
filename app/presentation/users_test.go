package presentation

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"production-go/data"
	"production-go/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestServer() (*http.ServeMux, *data.PostgreSQL) {
	// Initialize PostgreSQL connection
	db, err := data.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create service and setup routes
	serviceContainer := service.NewServiceContainer(db)
	SetupRoutes(mux, serviceContainer)

	return mux, db
}

func TestUserRegistration(t *testing.T) {
	mux, db := setupTestServer()
	defer db.DB().Close()

	// Create a test server using the mux
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Define test user credentials
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}{
		Username: "testuser",
		Password: "testpassword123",
		Email:    "testuser@example.com",
	}

	// Convert struct to JSON
	userJSON, _ := json.Marshal(user)

	// Test user registration
	res, err := http.Post(ts.URL+"/register", "application/json", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Cleanup: Delete the test user
	_, err = db.DB().Exec("DELETE FROM users WHERE username = $1", user.Username)
	assert.NoError(t, err)
}

func TestUserLogin(t *testing.T) {
	mux, db := setupTestServer()
	defer db.DB().Close()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Register a test user
	registerTestUser(t, ts, "testloginuser", "loginpassword123", "testloginuser@example.com")

	// Define login credentials
	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testloginuser",
		Password: "loginpassword123",
	}

	credentialsJSON, _ := json.Marshal(credentials)

	// Test user login
	res, err := http.Post(ts.URL+"/login", "application/json", bytes.NewBuffer(credentialsJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Cleanup
	cleanupTestUser(t, db, "testloginuser")
}

func TestChangePassword(t *testing.T) {
	mux, db := setupTestServer()
	defer db.DB().Close()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Register a test user
	registerTestUser(t, ts, "testchangepassworduser", "oldpassword123", "testchangepassworduser@example.com")

	// Define change password request
	changeRequest := struct {
		Username    string `json:"username"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}{
		Username:    "testchangepassworduser",
		OldPassword: "oldpassword123",
		NewPassword: "newpassword123",
	}

	changeRequestJSON, _ := json.Marshal(changeRequest)

	// Send change password request
	req, _ := http.NewRequest("POST", ts.URL+"/change-password", bytes.NewBuffer(changeRequestJSON))
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Verify login with new password
	loginWithNewPassword := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testchangepassworduser",
		Password: "newpassword123",
	}
	loginWithNewPasswordJSON, _ := json.Marshal(loginWithNewPassword)
	loginRes, err := http.Post(ts.URL+"/login", "application/json", bytes.NewBuffer(loginWithNewPasswordJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginRes.StatusCode)

	// Cleanup
	cleanupTestUser(t, db, "testchangepassworduser")
}

func registerTestUser(t *testing.T, ts *httptest.Server, username, password, email string) {
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}{
		Username: username,
		Password: password,
		Email:    email,
	}
	userJSON, _ := json.Marshal(user)
	_, err := http.Post(ts.URL+"/register", "application/json", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
}

func cleanupTestUser(t *testing.T, db *data.PostgreSQL, username string) {
	_, err := db.DB().Exec("DELETE FROM users WHERE username = $1", username)
	assert.NoError(t, err)
}
