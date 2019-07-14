package main

import (
	"cloud"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	CREATE_CLUSTER  = "create-cluster"
	DESTROY_CLUSTER = "destroy-cluster"
)

func printDefaultUsage() {
	fmt.Printf("usage: allspark <%s|%s>\n", CREATE_CLUSTER, DESTROY_CLUSTER)
}

func handleCreateCluster(cloudEnvironment string, template string) {
	start := time.Now().Second()

	client := cloud.Create(cloudEnvironment)
	client.CreateCluster(template)

	end := time.Now().Second()

	log.Printf("cluster is online after %d seconds\n", (end - start))
}

func handleDestroyCluster(clusterID string, cloudEnvironment string) {
	client := cloud.Create(cloudEnvironment)
	client.DestroyCluster(clusterID)
}

func main() {
	createCluster := flag.NewFlagSet("create-cluster", flag.ExitOnError)
	createCloudEnvironment := createCluster.String("cloud-environment", "",
		"Cloud environment; options include docker, aws, azure")
	template := createCluster.String("template", "",
		"/path/to/deployment-template")

	destroyCluster := flag.NewFlagSet("destroy-cluster", flag.ExitOnError)
	clusterID := destroyCluster.String("cluster-id", "",
		"ID of stack to be destroyed")
	destroyCloudEnvironment := destroyCluster.String("cloud-environment", "",
		"Cloud environment; options include docker, aws, azure")

	if len(os.Args) <= 1 {
		printDefaultUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case CREATE_CLUSTER:
		createCluster.Parse(os.Args[2:])
		if len(*createCloudEnvironment) == 0 ||
			len(*template) == 0 {
			createCluster.Usage()
			os.Exit(1)
		}
		handleCreateCluster(*createCloudEnvironment, *template)
	case DESTROY_CLUSTER:
		destroyCluster.Parse(os.Args[2:])
		if len(*clusterID) == 0 ||
			len(*destroyCloudEnvironment) == 0 {
			destroyCluster.Usage()
			os.Exit(1)
		}
		handleDestroyCluster(*clusterID, *destroyCloudEnvironment)
	default:
		printDefaultUsage()
		os.Exit(1)
	}
}
