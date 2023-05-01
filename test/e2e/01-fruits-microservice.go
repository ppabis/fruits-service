package e2e

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"

	dt "github.com/ory/dockertest/v3"
	d "github.com/ory/dockertest/v3/docker"
)

func FruitsMicroservice(redisHost string) (*dt.Resource, int, error) {
	port := 56081 + rand.Intn(200)

	err := setup()
	if err != nil {
		return nil, -1, err
	}

	res, err := pool.RunWithOptions(
		&dt.RunOptions{
			Repository:   "fruits-microservice",
			Tag:          "latest",
			ExposedPorts: []string{"8081/tcp"},
			PortBindings: map[d.Port][]d.PortBinding{
				"8081/tcp": {
					{
						HostPort: strconv.Itoa(port),
					},
				},
			},
			Env: []string{
				"REDIS_ENDPOINT=" + redisHost + ":6379",
				"PUBLIC_KEY_FILE=/run/secrets/public_key",
			},
			Mounts: []string{
				os.Getenv("PUBLIC_KEY_FILE") + ":/run/secrets/public_key",
			},
			NetworkID: network.Network.ID,
		},
		func(config *d.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = d.RestartPolicy{Name: "no"}
		})

	if err != nil {
		return nil, -1, err
	}

	if err := pool.Retry(func() error {
		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/")
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return err
		}
		return nil
	}); err != nil {
		return nil, -1, err
	}

	return res, port, err
}
