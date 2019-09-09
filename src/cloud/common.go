package cloud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"spark_monitor"
	"time"
	"util/serializer"
)

// Supported cloud environments
const (
	Aws    = "aws"
	Azure  = "azure"
	Docker = "docker"
)

const (
	masterIdentifier = "-master"
	workerIdentifier = "-worker-"
	sparkPort        = 7077
	aliveWorkers     = "Alive Workers:"
)

// CloudEnvironment base interface
type CloudEnvironment interface {
	CreateCluster() (string, error)
	DestroyCluster() error
	getClusterNodes() ([]string, error)
}

func waitForCluster(sparkWebURL string, expectedWorkerCount int,
	retryAttempts int) error {

	for i := 0; i < retryAttempts; i++ {
		workerCount, _ := getAliveWorkerCount(sparkWebURL)
		if workerCount == expectedWorkerCount {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("Spark cluster failed to launch")
}

func getAliveWorkerCount(sparkWebURL string) (int, error) {
	resp, err := http.Get(sparkWebURL + "/json/")
	if err == nil {
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var sparkClusterStatus spark_monitor.SparkClusterStatus

		err = serializer.Deserialize(contents, &sparkClusterStatus)
		if err != nil {
			return 0, err
		}

		return sparkClusterStatus.AliveWorkers, nil
	}
	return 0, err
}

// ReadTemplateConfiguration - reads a template configuration file
func ReadTemplateConfiguration(path string) ([]byte, error) {
	template, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer template.Close()

	return ioutil.ReadAll(template)
}

// Create a cloud environment (e.g. AWS, Docker, Azure, etc..)
func Create(environment string, clusterConfiguration []byte) (CloudEnvironment, error) {
	switch environment {
	case Aws:
		var result AwsEnvironment
		err := json.Unmarshal(clusterConfiguration, &result)
		return &result, err
	case Azure:
		return &AzureEnvironment{}, nil
	case Docker:
		var result DockerEnvironment
		err := json.Unmarshal(clusterConfiguration, &result)
		return &result, err
	default:
		log.Fatal("invalid cloud-environment " + environment)
	}

	return nil, nil
}
