package users

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func CreateUser(username, password string) error {
	// Creates a new user

	if len(password) < 8 {
		return fmt.Errorf("Password must be at least 8 characters")
	}

	passwordHash := hashPassword(password)
	if passwordHash == "" {
		return fmt.Errorf("Failed to hash password")
	}

	db, err := sql.Open("sqlite3", "monolith.db")
	defer db.Close()
	if err != nil {
		return err
	}

	if !ensureUsersTable(db) {
		return fmt.Errorf("Failed to ensure users table")
	}

	rows, err := db.Query("SELECT id FROM users WHERE username = ?", username)
	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("Username already exists")
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, passwordHash)

	return err
}
