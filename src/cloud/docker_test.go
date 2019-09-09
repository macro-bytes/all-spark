package cloud

import (
	"strconv"
	"testing"
	"util/serializer"
)

const (
	dockerTemplatePath = "../../sample_templates/docker.json"
)

func getDockerClient(t *testing.T) CloudEnvironment {
	templateConfig, err := ReadTemplateConfiguration(dockerTemplatePath)
	if err != nil {
		t.Fatal(err)
	}

	cloud, err := Create(Docker, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cloud
}

func TestCreateDockerCluster(t *testing.T) {
	cloud := getDockerClient(t)
	var spec DockerEnvironment

	err := serializer.DeserializePath(dockerTemplatePath, &spec)
	if err != nil {
		t.Fatal(err)
	}

	webURL, err := cloud.CreateCluster()
	if err != nil {
		t.Fatal(err)
	}

	clusterNodes, err := cloud.getClusterNodes()
	if err != nil {
		t.Fatal(err)
	}

	expectedNodeCount := spec.WorkerNodes + 1
	actualNodeCount := len(clusterNodes)

	if expectedNodeCount != actualNodeCount {
		t.Error("- expected " + strconv.Itoa(expectedNodeCount) + " spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}

	err = waitForCluster(webURL, spec.WorkerNodes, 20)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDestroyDockerCluster(t *testing.T) {
	cloud := getDockerClient(t)
	cloud.DestroyCluster()

	clusterNodes, err := cloud.getClusterNodes()
	if err != nil {
		t.Fatal(err)
	}

	actualNodeCount := len(clusterNodes)

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}
}
