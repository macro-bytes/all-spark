package cloud

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
	CreateCluster(configPath string) error
	DestroyCluster(identifier string) error
	getClusterNodes(identifier string) ([]string, error)
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
