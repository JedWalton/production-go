package service

import (
	"database/sql"
	"errors"
	"production-go/data"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	pg *data.PostgreSQL
}

func NewUserService(pg *data.PostgreSQL) *UserService {
	return &UserService{pg: pg}
}

func (s *UserService) RegisterUser(username, password, email string) error {
	// Validate input
	if err := validateRegistrationInput(username, password, email); err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user in the database
	return data.CreateUser(s.pg.DB(), username, string(hashedPassword), email)
}

func validateRegistrationInput(username, password, email string) error {
	// Check username (for example, ensure it's non-empty)
	if username == "" {
		return errors.New("username cannot be empty")
	}

	// Check password complexity
	// You can adjust the regex based on your password policy
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check email format
	if match, _ := regexp.MatchString(`^\S+@\S+\.\S+$`, email); !match {
		return errors.New("invalid email format")
	}

	return nil
}

func (s *UserService) LoginUser(username, password string) (bool, error) {
	// Retrieve the stored password hash from the data layer
	passwordHash, err := data.GetPasswordHash(s.pg.DB(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			return false, nil
		}
		// Other error
		return false, err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		// Password does not match
		return false, nil
	}

	// Login successful
	return true, nil
}
