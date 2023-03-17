package users

import (
	"database/sql"
	"fmt"
)

var sessions map[string]int = make(map[string]int)

func Authenticate(username, password string) (string, error) {
	// Authenticates user and returns a session cookie
	db, err := sql.Open("sqlite3", "monolith.db")
	defer db.Close()

	if err != nil {
		return "", err
	}

	if !ensureUsersTable(db) {
		return "", fmt.Errorf("Failed to ensure users table")
	}

	passwordHash := hashPassword(password)
	if passwordHash == "" {
		return "", fmt.Errorf("Failed to hash password")
	}

	user := getUser(db, username, passwordHash)
	user.Next()
	var id int
	err = user.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("Username or password not ok")
	}

	cookie := newCookie(username)
	sessions[cookie] = id

	return cookie, nil
}

func GetByCookie(cookie string) (int, error) {
	if id, ok := sessions[cookie]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("Invalid cookie")
}
