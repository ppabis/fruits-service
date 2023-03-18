package fruits

import "database/sql"

func ensureFruitsTable(db *sql.DB) bool {
	// Ensures the fruits table exists
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS fruits (id INTEGER PRIMARY KEY, user INTEGER, fruit TEXT)") // Could be foreign key
	return err == nil
}
