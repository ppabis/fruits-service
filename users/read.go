package users

import (
	"database/sql"
	"fmt"
	"monolith/config"
)

var sessions map[string]int = make(map[string]int)

func Authenticate(username, password string) (string, error) {
	// Authenticates user and returns a session cookie
	db, err := sql.Open("sqlite3", config.DbFile)
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

func Logout(cookie string) {
	// Logs out user
	delete(sessions, cookie)
}

func GetByCookie(cookie string) (int, error) {
	if id, ok := sessions[cookie]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("invalid cookie")
}

func GetUsername(db *sql.DB, id int) (string, error) {
	rows, err := db.Query("SELECT username FROM users WHERE id = ?", id)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("couldn't find user %d", id)
	}

	var username string
	err = rows.Scan(&username)

	if err != nil {
		return "", err
	}

	return username, nil
}
