package fruits

import (
	"database/sql"
	"strings"
)

func ensureFruitsTable(db *sql.DB) bool {
	// Ensures the fruits table exists
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS fruits (id INTEGER PRIMARY KEY, user INTEGER, fruit TEXT)") // Could be foreign key
	return err == nil
}

func isFruitSpecial(fruit string) bool {
	// Checks if a fruit is special
	fruit = strings.ToLower(fruit)
	return fruit == "pineapple" || fruit == "kiwi"
}

func hasCurrent(db *sql.DB, id int) (bool, error) {
	// Checks if a user has a current fruit
	rows, err := db.Query("SELECT id FROM fruits WHERE user = ?", id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
