package config

import "os"

var DbFile = "monolith.db"

func init() {
	dbFile := os.Getenv("USE_DB_FILE")
	if dbFile != "" {
		DbFile = dbFile
	}
}
