package cloud

type AzureEnvironment struct{}

func (e *AzureEnvironment) CreateCluster(configPath string) error {
	return nil
}

func (e *AzureEnvironment) DestroyCluster(identifier string) error {
	return nil
}

func (e *AzureEnvironment) getClusterNodes(identifier string) ([]string, error) {
	return []string{}, nil
}
