package cloud

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
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

// Environment base interface
type Environment interface {
	CreateCluster(templatePath string) (string, error)
	DestroyCluster(templatePath string) error
	getClusterNodes(identifier string) ([]string, error)
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
	resp, err := http.Get(sparkWebURL)
	if err == nil {
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		for _, el := range strings.Split(strip.StripTags(string(contents)),
			"\n") {

			buff := strings.TrimSpace(el)
			if strings.Contains(buff, aliveWorkers) {
				workerLine := strings.Split(buff, ": ")
				count, err := strconv.Atoi(strings.TrimSpace(workerLine[1]))
				if err != nil {
					return 0, err
				}
				return count, nil
			}
		}
	}
	return 0, err
}

func buildBaseIdentifier(identifier string) string {
	return identifier + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}

// Create a cloud environment (e.g. AWS, Docker, Azure, etc..)
func Create(environment string) Environment {
	switch environment {
	case Aws:
		return &AwsEnvironment{}
	case Azure:
		return &AzureEnvironment{}
	case Docker:
		return &DockerEnvironment{}
	default:
		log.Fatal("invalid cloud-environment " + environment)
	}

	return nil
}
