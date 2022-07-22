package daemon

import (
	"allspark/logger"
	"allspark/util/serializer"
	"os"
	"strings"
)

// AllSparkConfig - allspark configuration parameters struct
type AllSparkConfig struct {
	RedisHost                    string
	RedisPassword                string
	ClusterPendingTimeout        int64
	ClusterIdleTimeout           int64
	DoneReportTime               int64
	ClusterMaxRuntime            int64
	ClusterMaxTimeWithoutCheckin int64
	CancelTerminationDelay       int64
	AzureEnabled                 bool
	AwsEnabled                   bool
	DockerEnabled                bool
	CallbackURL                  string
}

var config AllSparkConfig

// Init - loads allspark configuration parameters into configParams
func Init(path string) {
	err := serializer.DeserializePath(path, &config)
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}
	redis_host := os.Getenv("ALLSPARK_REDIS_HOST")
	if len(redis_host) > 0 {
		if strings.Contains(redis_host, ":") {
			config.RedisHost = redis_host
		} else {
			logger.GetError().Printf("Skipping ALLSPARK_REDIS_HOST=%s as it is not host/port format", redis_host);
		}
	}
	logger.GetInfo().Printf("config parameters: %+v", config)
}

// GetAllSparkConfig - returns allspark configuration parameters
func GetAllSparkConfig() AllSparkConfig {
	return config
}
