package fruits

import (
	"database/sql"
	"fmt"
	"log"
	"monolith/config"
	"monolith/users"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
		log.Default().Printf("failed to ensure fruits table\n")
		return fmt.Errorf("failed to ensure fruits table")
	}

	if isFruitSpecial(name) && !users.IsUserSuper(id) {
		log.Default().Printf("%d is not allowed to set fruit %q\n", id, name)
		return fmt.Errorf("you are not allowed to have this fruit")
	}

	username, err := users.GetUsername(db, id)
	if err != nil {
		return err
	}

	err = setInSqlite(db, id, name)

	if err == nil {
		err = sendToFruitsMicroservice(id, username, name, users.IsUserSuper(id))
	} else {
		log.Default().Printf("setting fruit in SQLite failed: %v\n", err)
	}

	log.Default().Printf("setting fruit failed: %v\n", err)

	return err
}

func setInSqlite(db *sql.DB, user int, fruit string) error {
	do_update, err := hasCurrent(db, user)
	if err != nil {
		return err
	}

	if do_update {
		_, err = db.Exec("UPDATE fruits SET fruit = ? WHERE user = ?", fruit, user)
	} else {
		_, err = db.Exec("INSERT INTO fruits (user, fruit) VALUES (?, ?)", user, fruit)
	}

	return err
}

// Sends the PUT command to the fruits microservice to update the fruit
// on their end
func sendToFruitsMicroservice(id int, username string, fruit string, super bool) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	token, err := CreateTokenForFruits(id, username, super)
	if err != nil {
		log.Default().Printf("failed to create token for fruits microservice: %v\n", err)
		return err
	}

	form := url.Values{"fruit": []string{fruit}}
	req, err := http.NewRequest("PUT", config.FruitsEndpoint+"/fruit", strings.NewReader(form.Encode()))
	if err != nil {
		log.Default().Printf("failed to create request for fruits microservice: %v\n", err)
		return err
	}

	req.Header.Add("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(req)
	if err != nil {
		log.Default().Printf("request failed for fruits microservice: %v\n", err)
		return err
	}

	if resp.StatusCode >= 400 {
		body := make([]byte, 1024) // Read at most 1 kilobyte
		resp.Body.Read(body)
		resp.Body.Close()
		error_string := strings.Trim(string(body), "\x00")
		log.Default().Printf("fruits microservice returned %d: %s\n", resp.StatusCode, error_string)
		return fmt.Errorf("fruits microservice returned %d: %s", resp.StatusCode, error_string)
	}

	return nil
}
