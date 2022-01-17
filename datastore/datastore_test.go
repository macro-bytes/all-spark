package datastore

import (
	"allspark/daemon"
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

	daemon.Init("../daemon/allspark_config_test.json")
	client = GetRedisClient()

	if client == nil {
		t.Fatal("expected non-nil redis-client")
	}
}

func TestRunRedisCommand(t *testing.T) {
	daemon.Init("../daemon/allspark_config.json")
	client := GetRedisClient()
	if client == nil {
		t.Fatal("redis client object is nil")
	}

	success, err := client.HSet("test-map", "test-key", "test-value").Result()
	if err != nil {
		t.Fatal(err)
	}

	if !success {
		t.Fatal("Redis HSet returned non-successful status")
	}

	result, err := client.HDel("test-map", "test-key").Result()
	if err != nil {
		t.Fatal(err)
	}

	if result != 1 {
		t.Fatal("Redis HDel return non-zero status")
	}

	result, err = client.Del("test-map").Result()
	if err != nil {
		t.Fatal(err)
	}

	if result != 0 {
		t.Fatal("Redis HDel returned status greater than zero")
	}
}
