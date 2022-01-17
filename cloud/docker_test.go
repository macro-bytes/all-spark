package cloud

import (
	"allspark/util/serializer"
	"strconv"
	"testing"
)

const (
	dockerClusterTemplatePath    = "../dist/sample_templates/docker.json"
	dockerSingleNodeTemplatePath = "../dist/sample_templates/docker_single_node.json"
)

func getDockerClient(t *testing.T, templatePath string) CloudEnvironment {
	templateConfig, err := ReadTemplateConfiguration(templatePath)
	if err != nil {
		t.Fatal(err)
	}

	cloud, err := Create(Docker, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cloud
}

func TestCreateAndDestroyDockerCluster(t *testing.T) {
	cloud := getDockerClient(t, dockerClusterTemplatePath)
	var spec DockerEnvironment

	err := serializer.DeserializePath(dockerClusterTemplatePath, &spec)
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

	if cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned true; expected false")
	}
	cloud.DestroyCluster()

	if !cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned false; expected true")
	}

	clusterNodes, err = cloud.getClusterNodes()
	if err != nil {
		t.Fatal(err)
	}

	actualNodeCount = len(clusterNodes)

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}
}

func TestCreateAndDestroySingleNodeDockerCluster(t *testing.T) {
	cloud := getDockerClient(t, dockerSingleNodeTemplatePath)
	var spec DockerEnvironment

	err := serializer.DeserializePath(dockerSingleNodeTemplatePath, &spec)
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

	if cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned true; expected false")
	}
	cloud.DestroyCluster()

	if !cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned false; expected true")
	}

	clusterNodes, err = cloud.getClusterNodes()
	if err != nil {
		t.Fatal(err)
	}

	actualNodeCount = len(clusterNodes)

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}
}
