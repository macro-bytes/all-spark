// +build !cli

package main

import (
	"api"
	"cloud"
	"daemon"
	"logger"
	"monitor"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		logger.Fatal("usage: allspark-orchestrator /path/to/allspark_config.json")
		os.Exit(1)
	}

	daemon.Init(os.Args[1])
	go monitor.Run(-1,
		daemon.GetAllSparkConfig().ClusterMaxRuntime,
		daemon.GetAllSparkConfig().ClusterIdleTimeout,
		daemon.GetAllSparkConfig().ClusterPendingTimeout)

	switch daemon.GetAllSparkConfig().CloudEnvironment {
	case cloud.Aws:
		api.InitAwsAPI()
	case cloud.Docker:
		api.InitDockerAPI()
	default:
		logger.Fatal("invalid cloud environment specified")
	}
}
