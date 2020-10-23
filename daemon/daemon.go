package daemon

import (
	"allspark/logger"
	"allspark/util/serializer"
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
	logger.GetInfo().Printf("config parameters: %+v", config)
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}
}

// GetAllSparkConfig - returns allspark configuration parameters
func GetAllSparkConfig() AllSparkConfig {
	return config
}
