package presentation

import (
	"encoding/json"
	"net/http"
	"production-go/service"
)

func setupUsersServiceRoutes(mux *http.ServeMux, serviceContainer *service.ServiceContainer) {
	userService := serviceContainer.UserService

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Username string
			Password string
			Email    string
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := userService.RegisterUser(user.Username, user.Password, user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	})

	// Login user route
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Username string
			Password string
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		success, err := userService.LoginUser(user.Username, user.Password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !success {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	})

	// Change password route
	mux.HandleFunc("/change-password", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Username    string
			OldPassword string
			NewPassword string
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := userService.ChangePassword(user.Username, user.OldPassword, user.NewPassword); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password changed successfully"))
	})
}
