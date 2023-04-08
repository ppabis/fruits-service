package fruits

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"monolith/config"
	"monolith/users"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
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

	username, err := users.GetUsername(db, id)
	if err != nil {
		return err
	}

	if do_update {
		_, err = db.Exec("UPDATE fruits SET fruit = ? WHERE user = ?", name, id)
	} else {
		_, err = db.Exec("INSERT INTO fruits (user, fruit) VALUES (?, ?)", id, name)
	}

	if err == nil {
		err = setInRedis(id, username, name)
	}

	return err
}

func setInRedis(id int, username string, fruit string) error {
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisEndpoint,
	})
	ctx := context.TODO()
	if client == nil || client.Ping(ctx).Err() != nil {
		return fmt.Errorf("cannot connect to redis, aborting")
	}

	key := fmt.Sprintf("user:%d", id)
	username_b64 := base64.StdEncoding.EncodeToString([]byte(username))
	value := fmt.Sprintf("%s:%s", username_b64, fruit)

	return client.Set(ctx, key, value, 0).Err()
}
