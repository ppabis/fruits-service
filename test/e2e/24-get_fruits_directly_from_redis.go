package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

func GetFruitsDirectlyFromRedis(endpoint int) error {
	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:" + strconv.Itoa(endpoint),
	})

	ctx := context.TODO()
	res, err := redis.Get(ctx, "user:3").Result()
	if err != nil {
		return err
	}

	username_fruit := strings.Split(res, ":")

	if username_fruit[1] != "pineapple" {
		return fmt.Errorf("expected pineapple, got %s", username_fruit[1])
	}

	username, err := base64.StdEncoding.DecodeString(username_fruit[0])
	if err != nil {
		return err
	}

	if string(username) != "charlie" {
		return fmt.Errorf("expected charlie, got %s", username)
	}

	return nil
}
