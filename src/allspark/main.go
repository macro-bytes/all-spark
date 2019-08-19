// +build !lambda

package main

import (
	"cloud"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func printDefaultUsage() {
	fmt.Printf("usage: allspark <%s|%s>\n", CREATE_CLUSTER, DESTROY_CLUSTER)
}

func handleErrors(options *flag.FlagSet,
	cloudEnvironment string, templatePath string) {
	if len(cloudEnvironment) == 0 ||
		len(templatePath) == 0 {
		options.Usage()
		os.Exit(1)
	}
}

func handleCreateCluster(options *flag.FlagSet,
	cloudEnvironment string, templatePath string) {
	handleErrors(options, cloudEnvironment, templatePath)

	log.Printf("Launching spark cluster (note: this may take approximately 45 seconds).")
	start := time.Now().Unix()

	client := cloud.Create(cloudEnvironment)
	url, err := client.CreateCluster(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now().Unix()
	log.Printf("Cluster is online after %v seconds\n", (end - start))
	log.Printf("Master URL: " + url)
}

func handleDestroyCluster(options *flag.FlagSet,
	cloudEnvironment string, templatePath string) {
	handleErrors(options, cloudEnvironment, templatePath)
	client := cloud.Create(cloudEnvironment)
	client.DestroyCluster(templatePath)
}

func main() {
	createCluster := flag.NewFlagSet(CREATE_CLUSTER, flag.ExitOnError)
	createCloudEnvironment := createCluster.String("cloud-environment", "",
		"Cloud environment; options include docker, aws, azure")
	createTemplate := createCluster.String("template", "",
		"/path/to/deployment-template")

	destroyCluster := flag.NewFlagSet(DESTROY_CLUSTER, flag.ExitOnError)
	destroyCloudEnvironment := destroyCluster.String("cloud-environment", "",
		"Cloud environment; options include docker, aws, azure")
	destroyTemplate := destroyCluster.String("template", "",
		"/path/to/deployment-template")

	if len(os.Args) <= 1 {
		printDefaultUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case CREATE_CLUSTER:
		createCluster.Parse(os.Args[2:])
		handleCreateCluster(createCluster,
			*createCloudEnvironment, *createTemplate)
	case DESTROY_CLUSTER:
		destroyCluster.Parse(os.Args[2:])
		handleDestroyCluster(destroyCluster,
			*destroyCloudEnvironment, *destroyTemplate)
	default:
		printDefaultUsage()
		os.Exit(1)
	}
}
