package cloud

// AzureEnvironment interface
type AzureEnvironment struct{}

// CreateCluster - creates spark clusters
func (e *AzureEnvironment) CreateCluster() (string, error) {
	return "", nil
}

// DestroyCluster - destroys spark clusters
func (e *AzureEnvironment) DestroyCluster() error {
	return nil
}

func (e *AzureEnvironment) getClusterNodes() ([]string, error) {
	return []string{}, nil
}
