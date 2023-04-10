package fruits

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"monolith/config"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
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

	client := redis.NewClient(&redis.Options{
		Addr: config.RedisEndpoint,
	})
	ctx := context.TODO()

	fruits := make(map[string]string)
	for rows.Next() {
		var id int
		var fruit string
		var name string
		rows.Scan(&id, &fruit, &name)
		fruits[name] = fruit
		validateRecord(client, ctx, id, name, fruit)
	}

	return fruits, nil
}

func validateRecord(client *redis.Client, ctx context.Context, id int, username string, fruit string) {
	record, err := client.Get(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err != nil {
		log.Printf("ERROR: failed to get record for user %d: %v", id, err)
		return
	}

	splitRecord := strings.Split(record, ":")
	if len(splitRecord) != 2 {
		log.Printf("ERROR: invalid record: %s", record)
		return
	}

	decodedUsername, err := base64.StdEncoding.DecodeString(splitRecord[0])
	if err != nil {
		log.Printf("ERROR: failed to decode username: %v", err)
		return
	}

	if string(decodedUsername) != username || splitRecord[1] != fruit {
		log.Printf("ERROR: mismatched record for user %d: %s", id, record)
	}
}
