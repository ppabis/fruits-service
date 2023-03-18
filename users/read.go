package users

import (
	"database/sql"
	"fmt"
)

var sessions map[string]int = make(map[string]int)

func Authenticate(username, password string) (string, error) {
	// Authenticates user and returns a session cookie
	db, err := sql.Open("sqlite3", "monolith.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	if !ensureUsersTable(db) {
		return "", fmt.Errorf("failed to ensure users table")
	}

	passwordHash := hashPassword(password)
	if passwordHash == "" {
		return "", fmt.Errorf("failed to hash password")
	}

	user := getUser(db, username, passwordHash)
	if user == nil {
		return "", fmt.Errorf("username or password not ok")
	}
	defer user.Close()
	user.Next()
	var id int
	err = user.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("username or password not ok")
	}

	cookie := newCookie(username)
	sessions[cookie] = id

	return cookie, nil
}

func Logout(cookie string) error {
	// Logs out user
	delete(sessions, cookie)
	return nil
}

func GetByCookie(cookie string) (int, error) {
	if id, ok := sessions[cookie]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("invalid cookie")
}
