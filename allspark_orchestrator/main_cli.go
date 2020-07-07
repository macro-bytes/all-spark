// +build cli

package main

import (
	"allspark/cloud"
	"allspark/logger"
	"flag"
	"fmt"
	"os"
	"time"
)

func printDefaultUsage() {
	fmt.Printf("usage: allspark <%s|%s>\n", CreateCluster, DestroyCluster)
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

	logger.GetInfo().Println("Launching spark cluster (note: this may take up to 90 seconds).")
	start := time.Now().Unix()

	templateConfig, err := cloud.ReadTemplateConfiguration(templatePath)
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}

	client, err := cloud.Create(cloudEnvironment, templateConfig)
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}

	_, err = client.CreateCluster()
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}

	end := time.Now().Unix()
	logger.GetInfo().Printf("Cluster is online after %v seconds\n", (end - start))
}

func handleDestroyCluster(options *flag.FlagSet,
	cloudEnvironment string, templatePath string) {
	handleErrors(options, cloudEnvironment, templatePath)

	templateConfig, err := cloud.ReadTemplateConfiguration(templatePath)
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}

	client, err := cloud.Create(cloudEnvironment, templateConfig)
	if err != nil {
		logger.GetFatal().Fatalln(err)
	}
	client.DestroyCluster()
}

func main() {
	createCluster := flag.NewFlagSet(CreateCluster, flag.ExitOnError)
	createCloudEnvironment := createCluster.String("cloud-environment", "",
		"Cloud environment; options include docker, aws, azure")
	createTemplate := createCluster.String("template", "",
		"/path/to/deployment-template")

	destroyCluster := flag.NewFlagSet(DestroyCluster, flag.ExitOnError)
	destroyCloudEnvironment := destroyCluster.String("cloud-environment", "",
		"Cloud environment; options include docker, aws, azure")
	destroyTemplate := destroyCluster.String("template", "",
		"/path/to/deployment-template")

	if len(os.Args) <= 1 {
		printDefaultUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case CreateCluster:
		createCluster.Parse(os.Args[2:])
		handleCreateCluster(createCluster,
			*createCloudEnvironment, *createTemplate)
	case DestroyCluster:
		destroyCluster.Parse(os.Args[2:])
		handleDestroyCluster(destroyCluster,
			*destroyCloudEnvironment, *destroyTemplate)
	default:
		printDefaultUsage()
		os.Exit(1)
	}
}
