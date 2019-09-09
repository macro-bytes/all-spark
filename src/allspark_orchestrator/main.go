// +build !cli

package main

import (
	"api"
	"cloud"
	"daemon"
	"log"
	"monitor"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Panic("usage: allspark-orchestrator /path/to/allspark_config.json")
		os.Exit(1)
	}

	daemon.Init(os.Args[1])
	go monitor.MonitorSparkClusters(-1,
		daemon.GetAllSparkConfig().ClusterMaxRuntime,
		daemon.GetAllSparkConfig().ClusterMaxRuntime,
		daemon.GetAllSparkConfig().ClusterIdleTimeout)

	switch daemon.GetAllSparkConfig().CloudEnvironment {
	case cloud.Aws:
		api.InitAwsAPI()
	case cloud.Docker:
		api.InitDockerAPI()
	default:
		log.Panic("invalid cloud environment specified")
	}
}
