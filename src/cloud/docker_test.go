package cloud

import (
	"strconv"
	"template"
	"testing"
)

func TestCreateDockerCluster(t *testing.T) {
	var template template.DockerTemplate
	DeserializeTemplate("../../sample_templates/docker.json",
		&template)
	cloud := Create(Docker)

	webURL, err := cloud.CreateCluster("../../sample_templates/docker.json")
	if err != nil {
		t.Fatal(err)
	}

	clusterNodes, err := cloud.getClusterNodes(template.ClusterID)
	if err != nil {
		t.Fatal(err)
	}

	expectedNodeCount := template.WorkerNodes + 1
	actualNodeCount := len(clusterNodes)

	if expectedNodeCount != actualNodeCount {
		t.Error("- expected " + strconv.Itoa(expectedNodeCount) + " spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}

	err = waitForCluster(webURL, template.WorkerNodes, 20)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDestroyDockerCluster(t *testing.T) {
	templatePath := "../../sample_templates/docker.json"
	var template template.DockerTemplate
	DeserializeTemplate(templatePath, &template)

	cloud := Create(Docker)
	cloud.DestroyCluster(templatePath)

	clusterNodes, err := cloud.getClusterNodes(template.ClusterID)
	if err != nil {
		t.Fatal(err)
	}

	actualNodeCount := len(clusterNodes)

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}
}

func TestDeserializeTemplate(t *testing.T) {
	var template template.DockerTemplate

	err := DeserializeTemplate("does-not-exist", &template)
	if err == nil {
		t.Error("Expected non-nil error")
	}

	err = DeserializeTemplate("../../sample_templates/docker.json", &template)
	if err != nil {
		t.Error(err)
	}
}
