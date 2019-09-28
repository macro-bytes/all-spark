package cloud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"logger"
	"net/http"
	"os"
	"time"
	"util/serializer"
)

// SparkWorker describes the spark worker node state
type SparkWorker struct {
	ID            string `json:"id"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	WebUIAddress  string `json:"webuiaddress"`
	Cores         int    `json:"cores"`
	CoresUsed     int    `json:"coresused"`
	CoresFree     int    `json:"coresfree"`
	Memory        uint64 `json:"memory"`
	MemoryUsed    uint64 `json:"memoryused"`
	MemoryFree    uint64 `json:"memoryfree"`
	State         string `json:"state"`
	LastHeartBeat uint64 `json:"lastheartbeat"`
}

// SparkApp describes the spark application state
type SparkApp struct {
	ID             string `json:"id"`
	StartTime      uint64 `json:"starttime"`
	Name           string `json:"name"`
	Cores          int    `json:"cores"`
	User           string `json:"user"`
	MemoryPerSlave int    `json:"memoryperslave"`
	SubmitDate     string `json:"submitdate"`
	State          string `json:"state"`
	Duration       uint64 `json:"duration"`
}

// SparkClusterStatus describes the entire spark cluster state
type SparkClusterStatus struct {
	URL           string        `json:"url"`
	Workers       []SparkWorker `json:"workers"`
	AliveWorkers  int           `json:"aliveworkers"`
	Cores         int           `json:"cores"`
	CoresUsed     int           `json:"coresused"`
	Memory        uint64        `json:"memory"`
	MemoryUsed    uint64        `json:"memoryused"`
	ActiveApps    []SparkApp    `json:"activeapps"`
	CompletedApps []SparkApp    `json:"completedapps"`
	Status        string        `json:"status"`
}

// SparkStatusCheckIn - form body for the /checkin endpoint
type SparkStatusCheckIn struct {
	Status    SparkClusterStatus
	ClusterID string
}

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

		var sparkClusterStatus SparkClusterStatus

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
		logger.GetFatal().Fatalln(err)
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
	}

	return nil, errors.New("invalid cloud-environment " + environment)
}
