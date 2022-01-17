package cloud

import (
	"allspark/util/serializer"
	"strconv"
	"testing"
	"time"
)

const (
	azureClusterTemplatePath    = "../dist/sample_templates/azure.json"
	azureSingleNodeTemplatePath = "../dist/sample_templates/azure_single_node.json"
)

func getClient(t *testing.T, templatePath string) CloudEnvironment {
	templateConfig, err := ReadTemplateConfiguration(templatePath)
	if err != nil {
		t.Fatal(err)
	}

	cloud, err := Create(Azure, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cloud
}

func TestCreateListDeleteDisks(t *testing.T) {
	client1 := getClient(t, azureClusterTemplatePath).(*AzureEnvironment)
	client1.ClusterID = "azure-cluster-1"

	client2 := getClient(t, azureClusterTemplatePath).(*AzureEnvironment)
	client2.ClusterID = "azure-cluster-2"

	client1.createDisk(client1.ClusterID)
	client2.createDisk(client2.ClusterID)

	items, err := client1.getDisks()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 1 {
		t.Error("Expected array length of 1, got " + strconv.Itoa(len(items)))
	}

	if items[0] != client1.ClusterID {
		t.Error("Expected entry " + client1.ClusterID + ", got " + items[0])
	}

	client1.deleteDisk(client1.ClusterID)
	items, err = client1.getDisks()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 0 {
		t.Error("Expected array length of 0, got " + strconv.Itoa(len(items)))
	}

	items, err = client2.getDisks()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 1 {
		t.Error("Expected array length of 1, got " + strconv.Itoa(len(items)))
	}

	client2.deleteDisk(client2.ClusterID)
	items, err = client1.getDisks()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 0 {
		t.Error("Expected array length of 0, got " + strconv.Itoa(len(items)))
	}
}

func TestCreateListDeleteNics(t *testing.T) {
	client1 := getClient(t, azureClusterTemplatePath).(*AzureEnvironment)
	client1.ClusterID = "azure-cluster-1"

	client2 := getClient(t, azureClusterTemplatePath).(*AzureEnvironment)
	client2.ClusterID = "azure-cluster-2"

	client1.createNIC(client1.ClusterID)
	client2.createNIC(client2.ClusterID)

	items, err := client1.getNics()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 1 {
		t.Error("Expected array length of 1, got " + strconv.Itoa(len(items)))
	}

	if items[0] != client1.ClusterID {
		t.Error("Expected entry " + client1.ClusterID + ", got " + items[0])
	}

	client1.deleteNIC(client1.ClusterID)
	items, err = client1.getNics()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 0 {
		t.Error("Expected array length of 0, got " + strconv.Itoa(len(items)))
	}

	items, err = client2.getNics()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 1 {
		t.Error("Expected array length of 1, got " + strconv.Itoa(len(items)))
	}

	client2.deleteNIC(client2.ClusterID)
	items, err = client1.getNics()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 0 {
		t.Error("Expected array length of 0, got " + strconv.Itoa(len(items)))
	}
}

func TestGetPrimaryStorageKey(t *testing.T) {
	client := getClient(t, azureClusterTemplatePath)
	azureClient := client.(*AzureEnvironment)
	primaryKey, err := azureClient.getPrimaryStorageKey()
	if err != nil {
		t.Fatal(err)
	}

	if len(primaryKey) == 0 {
		t.Fatal("unable to retrieve primary storage account key")
	}
}

func TestCreateAndDestroyAzureCluster(t *testing.T) {
	cloud := getClient(t, azureClusterTemplatePath)
	var spec AzureEnvironment

	err := serializer.DeserializePath(azureClusterTemplatePath, &spec)
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

	if cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned true; expected false")
	}

	err = cloud.DestroyCluster()
	if err != nil {
		t.Fatal(err)
	}

	if cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned true; expected false")
	}

	time.Sleep(5 * time.Minute)

	if !cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned false; expected true")
	}

	clusterNodes, err = cloud.getClusterNodes()
	if err != nil {
		t.Error(err)
	}

	actualNodeCount = int64(len(clusterNodes))

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(int(actualNodeCount)) + " spark nodes.")
	}
}

func TestCreateAndDestroySingleNodeAzureCluster(t *testing.T) {
	cloud := getClient(t, azureSingleNodeTemplatePath)
	var spec AzureEnvironment

	err := serializer.DeserializePath(azureSingleNodeTemplatePath, &spec)
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

	if cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned true; expected false")
	}

	err = cloud.DestroyCluster()
	if err != nil {
		t.Fatal(err)
	}

	if cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned true; expected false")
	}

	time.Sleep(5 * time.Minute)

	if !cloud.DestructionConfirmed() {
		t.Error("DestructionConfirmed returned false; expected true")
	}

	clusterNodes, err = cloud.getClusterNodes()
	if err != nil {
		t.Error(err)
	}

	actualNodeCount = int64(len(clusterNodes))

	if 0 != actualNodeCount {
		t.Error("- expected 0 spark nodes.")
		t.Error("- got " + strconv.Itoa(int(actualNodeCount)) + " spark nodes.")
	}
}
