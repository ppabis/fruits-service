package fruits

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func GetFruits() (map[string]string, error) {
	// Returns a map of user ids to their fruits
	db, err := sql.Open("sqlite3", "monolith.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if !ensureFruitsTable(db) {
		return nil, fmt.Errorf("failed to ensure fruits table")
	}

	rows, err := db.Query("SELECT f.fruit, u.username FROM fruits f JOIN users u ON f.user = u.id")
	if err != nil {
		return nil, err
	}

	fruits := make(map[string]string)
	for rows.Next() {
		var fruit string
		var name string
		rows.Scan(&fruit, &name)
		fruits[name] = fruit
	}

	return fruits, nil
}
