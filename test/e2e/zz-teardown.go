package e2e

import (
	dt "github.com/ory/dockertest/v3"
)

func Teardown(resources ...*dt.Resource) {
	for _, r := range resources {
		if r != nil {
			pool.Purge(r)
		}
	}

	if network != nil {
		pool.RemoveNetwork(network)
	}
}
