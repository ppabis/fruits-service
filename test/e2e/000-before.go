package e2e

import (
	"log"
	"math/rand"
	"monolith/config"
	"monolith/router"
	"net/http"
	"strconv"
	"testing"

	dt "github.com/ory/dockertest/v3"
)

var fruitsContainer *dt.Resource
var redisContainer *dt.Resource
var httpClient = &http.Client{}

// Starts redis container, fruits container
// the monolith server and creates test users
func Before(t *testing.T) *struct {
	FruitsHost   string
	RedisHost    string
	RedisPort    int
	MonolithHost string
} {
	monolithPort := 58080 + rand.Intn(200)
	var fruitsPort int = -1
	var redisHost string
	var redisPort = -1
	var err error

	redisContainer, redisHost, redisPort, err = Redis()
	if err != nil {
		t.Error(err)
		return nil
	}

	log.Default().Printf("Redis available at %s\n", redisHost)

	fruitsContainer, fruitsPort, err = FruitsMicroservice(redisHost)
	if err != nil {
		t.Error(err)
		return nil
	}

	log.Default().Printf("Fruits available at localhost:%d\n", fruitsPort)

	config.FruitsEndpoint = "http://localhost:" + strconv.Itoa(fruitsPort)

	err = CreateUsers()
	if err != nil {
		t.Error(err)
		return nil
	}

	go router.Serve(monolithPort)

	return &struct {
		FruitsHost   string
		RedisHost    string
		RedisPort    int
		MonolithHost string
	}{
		FruitsHost:   "http://localhost:" + strconv.Itoa(fruitsPort),
		RedisHost:    redisHost,
		RedisPort:    redisPort,
		MonolithHost: "http://localhost:" + strconv.Itoa(monolithPort),
	}
}
