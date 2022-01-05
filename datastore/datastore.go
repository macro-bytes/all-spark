package datastore

import (
	"allspark/daemon"
	"allspark/logger"

	"github.com/go-redis/redis"
)

// GetRedisClient - returns Redis client
func GetRedisClient() *redis.Client {
	config := daemon.GetAllSparkConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       0,
	})

	result, err := client.Ping().Result()
	if err != nil {
		logger.GetError().Println(err)
		defer client.Close()
	}

	if result != "PONG" {
		logger.GetError().Println("unable to connect to redis; server did not respond to ping")
		defer client.Close()
	}

	return client
}
