package cloud

import (
	"strconv"
	"testing"
	"util/serializer"
)

const (
	awsTemplatePath = "../../sample_templates/aws.json"
)

func getAwsClient(t *testing.T) CloudEnvironment {
	templateConfig, err := ReadTemplateConfiguration(awsTemplatePath)
	if err != nil {
		t.Fatal(err)
	}

	cloud, err := Create(Aws, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cloud
}

func TestCreateAwsCluster(t *testing.T) {
	cloud := getAwsClient(t)
	var spec AwsEnvironment

	err := serializer.DeserializePath(awsTemplatePath, &spec)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cloud.CreateCluster()
	if err != nil {
		t.Fatal(err)
	}

	clusterNodes, err := cloud.getClusterNodes()
	if err != nil {
		t.Error(err)
	}

	expectedNodeCount := spec.WorkerNodes + 1
	actualNodeCount := int64(len(clusterNodes))

	if expectedNodeCount != actualNodeCount {
		t.Error("- expected " + strconv.FormatInt(expectedNodeCount, 10) +
			" spark nodes.")
		t.Error("- got " + strconv.FormatInt(actualNodeCount, 10) +
			" spark nodes.")
	}
}

func TestDestroyAwsCluster(t *testing.T) {
	cloud := getAwsClient(t)
	cloud.DestroyCluster()

	clusterNodes, err := cloud.getClusterNodes()
	if err != nil {
		t.Error(err)
	}

	actualNodeCount := len(clusterNodes)

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}
}
