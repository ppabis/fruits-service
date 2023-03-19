package users

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"monolith/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newCookie(username string) string {
	// Generates a random string
	input := fmt.Sprintf("%d+%s", time.Now().UnixNano(), username)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

func ensureUsersTable(db *sql.DB) bool {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT, super INTEGER DEFAULT 0)")
	return err == nil
}

func IsUserSuper(id int) bool {
	// Checks if user is super
	db, err := sql.Open("sqlite3", config.DbFile)
	if err != nil {
		return false
	}
	defer db.Close()

	if !ensureUsersTable(db) {
		return false
	}
	rows, err := db.Query("SELECT super FROM users WHERE id = ?", id)
	if err != nil {
		return false
	}
	defer rows.Close()

	rows.Next()
	var super int
	err = rows.Scan(&super)
	if err != nil {
		return false
	}

	return super > 0
}

func hashPassword(password string) string {
	// Hashes password using bcrypt
	const static_salt = "monolith"
	bytes := sha256.Sum256([]byte(password + static_salt))

	return base64.URLEncoding.EncodeToString(bytes[:])

}

func getUser(db *sql.DB, username string, passwordHash string) *sql.Rows {
	user, err := db.Query("SELECT id FROM users WHERE username = ? AND password = ?", username, passwordHash)
	if err != nil {
		return nil
	}
	return user
}
