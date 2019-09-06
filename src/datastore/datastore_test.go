package datastore

import (
	"allspark_config"
	"testing"
)

func TestGetRedisClient(t *testing.T) {
	allspark_config.Init("../allspark_config/allspark_config.json")
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
