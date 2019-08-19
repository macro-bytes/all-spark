package cloud

// AzureEnvironment interface
type AzureEnvironment struct{}

// CreateCluster - creates spark clusters
func (e *AzureEnvironment) CreateCluster(configPath string) (string, error) {
	return "", nil
}

// DestroyCluster - destroys spark clusters
func (e *AzureEnvironment) DestroyCluster(identifier string) error {
	return nil
}

func (e *AzureEnvironment) getClusterNodes(identifier string) ([]string, error) {
	return []string{}, nil
}
