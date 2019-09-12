// +build !cli

package main

import (
	"api"
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

	api.Init()
}
