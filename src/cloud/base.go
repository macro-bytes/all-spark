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

const (
	AWS               = "aws"
	AZURE             = "azure"
	DOCKER            = "docker"
	MASTER_IDENTIFIER = "-master"
	WORKER_IDENTIFIER = "-worker-"
	SPARK_PORT        = 7077
	ALIVE_WORKERS     = "Alive Workers:"
)

type Environment interface {
	CreateCluster(templatePath string) (string, error)
	DestroyCluster(templatePath string) error
	getClusterNodes(identifier string) ([]string, error)
}

func waitForCluster(sparkWebUrl string, expectedWorkerCount int,
	retryAttempts int) error {

	for i := 0; i < retryAttempts; i++ {
		workerCount, _ := getAliveWorkerCount(sparkWebUrl)
		if workerCount == expectedWorkerCount {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("Spark cluster failed to launch")
}

func getAliveWorkerCount(sparkWebUrl string) (int, error) {
	resp, err := http.Get(sparkWebUrl)
	if err == nil {
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		for _, el := range strings.Split(strip.StripTags(string(contents)),
			"\n") {

			buff := strings.TrimSpace(el)
			if strings.Contains(buff, ALIVE_WORKERS) {
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

func Create(environment string) Environment {
	switch environment {
	case AWS:
		return &AwsEnvironment{}
	case AZURE:
		return &AzureEnvironment{}
	case DOCKER:
		return &DockerEnvironment{}
	default:
		log.Fatal("invalid cloud-environment " + environment)
	}

	return nil
}
