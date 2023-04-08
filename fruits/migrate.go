package fruits

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"monolith/config"

	_ "github.com/mattn/go-sqlite3"

	"github.com/redis/go-redis/v9"
)

func MigrateSQLiteToRedis() {
	// This function will migrate current SQL `fruits` table structure
	// to a Redis database.
	userFruitMap, err := getUserAndFruit()
	if err != nil {
		fmt.Printf("failed to get user and fruit table: %v\n", err)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr: config.RedisEndpoint,
	})

	ctx := context.TODO()
	if client == nil || client.Ping(ctx).Err() != nil {
		fmt.Printf("failed to create redis client with endpoint %s\n", config.RedisEndpoint)
		return
	}

	var counter = 0
	for id, userFruit := range userFruitMap {
		// Store username and fruit as a tuple formatted as
		// base64(<username>):<fruit> under key user:<id>
		key := fmt.Sprintf("user:%d", id)
		user_b64 := base64.StdEncoding.EncodeToString([]byte(userFruit.username))
		value := fmt.Sprintf("%s:%s", user_b64, userFruit.fruit)

		err := client.Set(ctx, key, value, 0).Err()
		if err != nil {
			fmt.Printf("failed to set key %s: %v\n", key, err)
		} else {
			counter++
		}
	}

	fmt.Printf("migrated %d records to redis\n", counter)
}

type userAndFruit struct {
	username string
	fruit    string
}

func getUserAndFruit() (map[int]userAndFruit, error) {
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

	userFruitMap := make(map[int]userAndFruit)
	for rows.Next() {
		var id int
		var fruit string
		var name string
		rows.Scan(&id, &fruit, &name)
		userFruitMap[id] = userAndFruit{name, fruit}
	}

	return userFruitMap, nil
}
