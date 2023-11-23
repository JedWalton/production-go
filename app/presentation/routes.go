package presentation

import (
	"net/http"
	"production-go/service"
)

// SetupRoutes initializes routes for the user service
func SetupRoutes(mux *http.ServeMux, serviceContainer *service.ServiceContainer) {
	// Register user route
	setupUsersServiceRoutes(mux, serviceContainer)
}
