package fruits

import (
	"database/sql"
	"fmt"
	"monolith/users"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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
	return rows.Next(), nil
}

func UpdateFruit(id int, name string) error {
	// Updates a fruit
	db, err := sql.Open("sqlite3", "monolith.db")
	defer db.Close()
	if err != nil {
		return err
	}

	if !ensureFruitsTable(db) {
		return fmt.Errorf("Failed to ensure fruits table")
	}

	if isFruitSpecial(name) && !users.IsUserSuper(id) {
		return fmt.Errorf("You are not allowed to have this fruit")
	}

	do_update, err := hasCurrent(db, id)
	if err != nil {
		return err
	}

	if do_update {
		_, err = db.Exec("UPDATE fruits SET name = ? WHERE user = ?", name, id)
	} else {
		_, err = db.Exec("INSERT INTO fruits (user, name) VALUES (?, ?)", id, name)
	}

	return err
}
