package main

import (
	"cloud"
	"flag"
	"fmt"
	"os"
)

const (
	CREATE_STACK  = "create-stack"
	DESTROY_STACK = "destroy-stack"
)

func printDefaultUsage() {
	fmt.Println("usage: allspark <create-stack|destroy-stack>")
}

func handleCreateCluster(cloudEnvironment string, template string) {
	client := cloud.Create(cloudEnvironment)
	client.CreateCluster(template)
}

func handleDestroyCluster(clusterID string) {
	client := cloud.Create("")
	client.DestroyCluster(clusterID)
}

func main() {
	createStack := flag.NewFlagSet("create-cluster", flag.ExitOnError)
	cloudEnvironment := createStack.String("cloud_environment", "aws",
		"Cloud environment; options include aws and azure")
	template := createStack.String("template", "",
		"/path/to/deployment-template")

	destroyStack := flag.NewFlagSet("destroy-cluster", flag.ExitOnError)
	clusterID := destroyStack.String("cluster-id", "",
		"ID of stack to be destroyed")

	if len(os.Args) <= 1 {
		printDefaultUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case CREATE_STACK:
		createStack.Parse(os.Args[2:])
		handleCreateCluster(*cloudEnvironment, *template)
	case DESTROY_STACK:
		destroyStack.Parse(os.Args[2:])
		handleDestroyCluster(*clusterID)
	default:
		printDefaultUsage()
		os.Exit(1)
	}
}
