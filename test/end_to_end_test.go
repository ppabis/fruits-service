package test

import (
	"monolith/test/e2e"
	"testing"
)

func Test_End_To_End(t *testing.T) {

	var err error

	conf := e2e.Before(t)
	defer e2e.After(t)
	if conf == nil {
		t.Error("Could not setup test")
		return
	}

	err = e2e.SetFruits(conf.MonolithHost)
	if err != nil {
		t.Error(err)
	}

	err = e2e.GetFruits(conf.MonolithHost)
	if err != nil {
		t.Error(err)
	}

	err = e2e.GetFruitsMicroservice(conf.MonolithHost)
	if err != nil {
		t.Error(err)
	}

	err = e2e.GetFruitsDirectlyFromMicroservice(conf.FruitsHost)
	if err != nil {
		t.Error(err)
	}

	err = e2e.GetFruitsDirectlyFromRedis(conf.RedisPort)
	if err != nil {
		t.Error(err)
	}

}
