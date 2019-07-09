package cloud

const (
	AWS   = "aws"
	AZURE = "azure"
)

type Environment interface {
	CreateCluster() (string, error)
	DestroyCluster() error
	save() error
}

func Create(environment string) Environment {
	switch environment {
	case AWS:
		return &AwsEnvironment{}
	case AZURE:
		return &AzureEnvironment{}
	default:
		panic("invalid cloud-environment " + environment)
	}
}
