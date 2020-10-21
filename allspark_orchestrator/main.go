// +build !cli

package main

import (
	"allspark/api"
	"allspark/daemon"
	"allspark/logger"
	"allspark/monitor"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		logger.GetFatal().Fatalln("usage: allspark-orchestrator /path/to/allspark_config.json")
		os.Exit(1)
	}

	daemon.Init(os.Args[1])
	go monitor.Run(-1,
		daemon.GetAllSparkConfig().ClusterMaxRuntime,
		daemon.GetAllSparkConfig().ClusterIdleTimeout,
		daemon.GetAllSparkConfig().ClusterMaxTimeWithoutCheckin,
		daemon.GetAllSparkConfig().ClusterPendingTimeout,
		daemon.GetAllSparkConfig().DoneReportTime,
		daemon.GetAllSparkConfig().CancelTerminationDelay)

	api.Init()
}
