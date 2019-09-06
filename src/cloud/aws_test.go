package cloud

import (
	"strconv"
	"template"
	"testing"
)

func TestCreateAwsCluster(t *testing.T) {
	var template template.AwsTemplate
	DeserializeTemplate("../../sample_templates/aws.json",
		&template)

	cloud := Create(Aws)
	_, err := cloud.CreateCluster("../../sample_templates/aws.json")
	if err != nil {
		t.Fatal(err)
	}

	clusterNodes, err := cloud.getClusterNodes(template.ClusterID)
	if err != nil {
		t.Error(err)
	}

	expectedNodeCount := template.WorkerNodes + 1
	actualNodeCount := int64(len(clusterNodes))

	if expectedNodeCount != actualNodeCount {
		t.Error("- expected " + strconv.FormatInt(expectedNodeCount, 10) +
			" spark nodes.")
		t.Error("- got " + strconv.FormatInt(actualNodeCount, 10) +
			" spark nodes.")
	}

	/*
		err = waitForCluster(webUrl, int(template.WorkerNodes), 10)
		if err != nil {
			t.Fatal(err)
		}
	*/
}

func TestDestroyAwsCluster(t *testing.T) {
	templatePath := "../../sample_templates/aws.json"
	var template template.AwsTemplate
	DeserializeTemplate(templatePath, &template)

	cloud := Create(Aws)
	cloud.DestroyCluster(templatePath)

	clusterNodes, err := cloud.getClusterNodes(template.ClusterID)
	if err != nil {
		t.Error(err)
	}

	actualNodeCount := len(clusterNodes)

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(actualNodeCount) + " spark nodes.")
	}
}
