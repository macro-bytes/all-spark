package cloud

import (
	"allspark/util/serializer"
	"strconv"
	"testing"
	"time"
)

const (
	azureTemplatePath = "../dist/sample_templates/azure.json"
)

func getClient(t *testing.T) CloudEnvironment {
	templateConfig, err := ReadTemplateConfiguration(azureTemplatePath)
	if err != nil {
		t.Fatal(err)
	}

	cloud, err := Create(Azure, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cloud
}

func TestGetDisks(t *testing.T) {
}

func TestGetNics(t *testing.T) {
}

func TestGetPrimaryStorageKey(t *testing.T) {
	client := getClient(t)
	azureClient := client.(*AzureEnvironment)
	primaryKey, err := azureClient.getPrimaryStorageKey()
	if err != nil {
		t.Fatal(err)
	}

	if len(primaryKey) == 0 {
		t.Fatal("unable to retrieve primary storage account key")
	}
}

func TestCreateAzureCluster(t *testing.T) {
	cloud := getClient(t)
	var spec AzureEnvironment

	err := serializer.DeserializePath(azureTemplatePath, &spec)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cloud.CreateCluster()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Minute)

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

func TestDestroyAzureCluster(t *testing.T) {
	cloud := getClient(t)
	var spec AzureEnvironment

	err := serializer.DeserializePath(azureTemplatePath, &spec)
	if err != nil {
		t.Fatal(err)
	}

	err = cloud.DestroyCluster()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Minute)
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
