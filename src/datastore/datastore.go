package datastore

import (
	"allspark_config"

	"github.com/go-redis/redis"
)

// GetRedisClient - returns Redis client
func GetRedisClient() *redis.Client {
	config := allspark_config.GetAllSparkConfig()
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       0,
	})
}
