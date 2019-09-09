// +build !cli

package main

import (
	"allspark_config"
	"api"
	"cloud"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Panic("usage: allspark-orchestrator /path/to/allspark_config.json")
		os.Exit(1)
	}

	allspark_config.Init(os.Args[1])

	switch allspark_config.GetAllSparkConfig().CloudEnvironment {
	case cloud.Aws:
		api.InitAwsAPI()
	case cloud.Docker:
		api.InitDockerAPI()
	default:
		log.Panic("invalid cloud environment specified")
	}
}
