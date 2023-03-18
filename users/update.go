package users

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func UpdateUserSuperStatus(id int, super int) error {
	// Updates user super status
	db, err := sql.Open("sqlite3", "monolith.db")
	if err != nil {
		return err
	}
	defer db.Close()

	if !ensureUsersTable(db) {
		return fmt.Errorf("failed to ensure users table")
	}
	_, err = db.Exec("UPDATE users SET super = ? WHERE id = ?", super, id)
	return err

}
