package datastore

import (
	"daemon"
	"testing"
)

func TestGetRedisClient(t *testing.T) {
	daemon.Init("../daemon/allspark_config.json")
	client := GetRedisClient()
	if client == nil {
		t.Fatal("redis client object is nil")
	}

	result, err := client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}

	if result != "PONG" {
		t.Fatal("unable to ping redis server")
	}
}
