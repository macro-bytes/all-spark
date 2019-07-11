package cloud

import (
	"container/list"
)

type AwsEnvironment struct {
	instanceIDs *list.List
}

func (e *AwsEnvironment) CreateCluster(configPath string) error {
	if e.instanceIDs == nil {
		e.instanceIDs = list.New()
	}

	return nil
}

func (e *AwsEnvironment) DestroyCluster(identifier string) error {
	return nil
}

func (e *AwsEnvironment) getClusterNodes(identifier string) ([]string, error) {
	return []string{}, nil
}
