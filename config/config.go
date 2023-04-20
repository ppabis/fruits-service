package config

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var DbFile = "monolith.db"
var RedisEndpoint = "localhost:6379"
var TokenLifetime = 30 * time.Second
var PrivateKey *ecdsa.PrivateKey

func init() {
	var err error = nil

	dbFile := os.Getenv("USE_DB_FILE")
	if dbFile != "" {
		DbFile = dbFile
	}

	redisEndpoint := os.Getenv("USE_REDIS_ENDPOINT")
	if redisEndpoint != "" {
		RedisEndpoint = redisEndpoint
	}

	privateKey := []byte(os.Getenv("PRIVATE_KEY"))
	if len(privateKey) == 0 {
		privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
		if privateKeyFile == "" {
			privateKeyFile = "server.pem"
		}
		privateKey, err = os.ReadFile(privateKeyFile)
		if err != nil {
			panic(fmt.Errorf("could not read private key: %w", err))
		}
	}

	PrivateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKey)
	if err != nil {
		panic(fmt.Errorf("could not parse private key: %w", err))
	}

}
