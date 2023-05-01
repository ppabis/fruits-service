package e2e

import (
	"context"
	"math/rand"
	"strconv"

	dt "github.com/ory/dockertest/v3"
	d "github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
)

var pool *dt.Pool
var network *dt.Network

// Prepapres the docker pool and network
func setup() error {
	var err error
	if pool == nil {
		pool, err = dt.NewPool("")
		if err != nil {
			return err
		}
	}

	if network == nil {
		network, err = pool.CreateNetwork("e2e")
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns resource for redis container, port,
// container name (host) and error
func Redis() (*dt.Resource, string, int, error) {
	port := 55037 + rand.Intn(200)
	name := "redis" + strconv.Itoa(port)

	err := setup()
	if err != nil {
		return nil, "", -1, err
	}

	res, err := pool.RunWithOptions(
		&dt.RunOptions{
			Repository:   "redis",
			Tag:          "latest",
			Name:         name,
			ExposedPorts: []string{"6379/tcp"},
			PortBindings: map[d.Port][]d.PortBinding{
				"6379/tcp": {
					{
						HostPort: strconv.Itoa(port),
					},
				},
			},
			NetworkID: network.Network.ID,
		},
		func(config *d.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = d.RestartPolicy{Name: "no"}
		})

	if err != nil {
		return nil, "", -1, err
	}

	if err := pool.Retry(func() error {
		return redis.
			NewClient(&redis.Options{
				Addr: "127.0.0.1:" + strconv.Itoa(port),
			}).
			Ping(context.TODO()).
			Err()
	}); err != nil {
		return nil, "", -1, err
	}

	return res, name, port, err
}
