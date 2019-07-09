package cloud

type AzureEnvironment struct{}

func (e *AzureEnvironment) CreateCluster() (string, error) {
	return "azure", nil
}

func (e *AzureEnvironment) DestroyCluster() error {
	return nil
}

func (e *AzureEnvironment) save() error {
	return nil
}
