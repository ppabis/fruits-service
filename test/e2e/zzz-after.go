package e2e

import (
	"monolith/router"
	"os"
	"testing"
)

func After(t *testing.T) {
	Teardown(fruitsContainer)
	Teardown(redisContainer)
	os.Remove(os.Getenv("USE_DB_FILE"))
	router.Shutdown()
}
