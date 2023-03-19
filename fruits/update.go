package fruits

import (
	"database/sql"
	"fmt"
	"monolith/config"
	"monolith/users"

	_ "github.com/mattn/go-sqlite3"
)

func UpdateFruit(id int, name string) error {
	// Updates a fruit
	db, err := sql.Open("sqlite3", config.DbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if !ensureFruitsTable(db) {
		return fmt.Errorf("failed to ensure fruits table")
	}

	if isFruitSpecial(name) && !users.IsUserSuper(id) {
		return fmt.Errorf("you are not allowed to have this fruit")
	}

	do_update, err := hasCurrent(db, id)
	if err != nil {
		return err
	}

	if do_update {
		_, err = db.Exec("UPDATE fruits SET fruit = ? WHERE user = ?", name, id)
	} else {
		_, err = db.Exec("INSERT INTO fruits (user, fruit) VALUES (?, ?)", id, name)
	}

	return err
}
