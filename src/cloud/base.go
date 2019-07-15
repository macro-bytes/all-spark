package cloud

import (
	"strconv"
	"time"
)

const (
	AWS               = "aws"
	AZURE             = "azure"
	DOCKER            = "docker"
	MASTER_IDENTIFIER = "-master"
	WORKER_IDENTIFIER = "-worker-"
	SPARK_BASE_IMAGE  = "mshoaazar/spark-2.4-standalone:latest"
	SPARK_PORT        = "7077"
)

type Environment interface {
	CreateCluster(templatePath string) error
	DestroyCluster(identifier string) error
	getClusterNodes(identifier string) ([]string, error)
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
		panic("invalid cloud-environment " + environment)
	}
}
