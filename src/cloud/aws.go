package cloud

import (
	"container/list"
)

type AwsEnvironment struct {
	instanceIDs *list.List
}

func (e *AwsEnvironment) CreateCluster() (string, error) {
	if e.instanceIDs == nil {
		e.instanceIDs = list.New()
	}

	return "aws", nil
}

func (e *AwsEnvironment) DestroyCluster() error {
	return nil
}

func (e *AwsEnvironment) save() error {
	return nil
}
