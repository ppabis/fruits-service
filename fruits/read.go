package fruits

import (
	"database/sql"
	"fmt"
	"monolith/config"

	_ "github.com/mattn/go-sqlite3"
)

func GetFruits() (map[string]string, error) {
	// Returns a map of user ids to their fruits
	db, err := sql.Open("sqlite3", config.DbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if !ensureFruitsTable(db) {
		return nil, fmt.Errorf("failed to ensure fruits table")
	}

	rows, err := db.Query("SELECT u.id, f.fruit, u.username FROM fruits f JOIN users u ON f.user = u.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fruits := make(map[string]string)
	for rows.Next() {
		var id int
		var fruit string
		var name string
		rows.Scan(&id, &fruit, &name)
		fruits[name] = fruit
	}

	return fruits, nil
}
