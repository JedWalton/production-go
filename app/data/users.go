package data

import (
	"database/sql"
	// other imports
)

func CreateUser(db *sql.DB, username, passwordHash, email string) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare SQL statement within the transaction
	stmt, err := tx.Prepare("INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback() // Important to rollback if there's an error
		return err
	}
	defer stmt.Close() // Ensure the statement is closed after execution

	// Execute the statement
	_, err = stmt.Exec(username, passwordHash, email)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func DeleteUser(db *sql.DB, username string) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the delete statement within the transaction
	stmt, err := tx.Prepare("DELETE FROM users WHERE username = $1")
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return err
	}
	defer stmt.Close() // Ensure the statement is closed after execution

	// Execute the statement
	_, err = stmt.Exec(username)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GetPasswordHash retrieves the password hash for a given username.
func GetPasswordHash(db *sql.DB, username string) (string, error) {
	var passwordHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func UpdateUserPassword(db *sql.DB, username, newHashedPassword string) error {
	// Prepare SQL statement to update user's password
	stmt, err := db.Prepare("UPDATE users SET password_hash = $1 WHERE username = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(newHashedPassword, username)
	return err
}
