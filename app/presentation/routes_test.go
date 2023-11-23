package presentation

import (
	"net/http"
	"net/http/httptest"
	"production-go/data"
	"production-go/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	// Initialize database connection
	db, err := data.NewPostgreSQL()
	assert.NoError(t, err)
	defer db.DB().Close()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create service and setup routes
	serviceContainer := service.NewServiceContainer(db)
	SetupRoutes(mux, serviceContainer)

	// Create a test server using the mux
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Define a list of expected routes and methods
	routes := []struct {
		path   string
		method string
	}{
		{"/register", "POST"},
		{"/login", "POST"},
		{"/change-password", "POST"},
		// Add other routes as necessary
	}

	for _, route := range routes {
		req, _ := http.NewRequest(route.method, ts.URL+route.path, nil)
		res, err := http.DefaultClient.Do(req)

		// Check that the route exists and is accessible
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusNotFound, res.StatusCode, "Route not found: %s", route.path)
	}
}
