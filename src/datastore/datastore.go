package datastore

import (
	"daemon"

	"github.com/go-redis/redis"
)

// GetRedisClient - returns Redis client
func GetRedisClient() *redis.Client {
	config := daemon.GetAllSparkConfig()
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       0,
	})
}
