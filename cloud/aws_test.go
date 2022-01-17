package cloud

import (
	"allspark/util/serializer"
	"regexp"
	"strconv"
	"testing"
	"time"
)

const (
	awsClusterTemplatePath    = "../dist/sample_templates/aws.json"
	awsSingleNodeTemplatePath = "../dist/sample_templates/aws_single_node.json"
)

func getAwsClient(t *testing.T, templatePath string) CloudEnvironment {
	templateConfig, err := ReadTemplateConfiguration(templatePath)
	if err != nil {
		t.Fatal(err)
	}

	cloud, err := Create(Aws, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cloud
}

func TestResolveAMI(t *testing.T) {
	var spec AwsEnvironment

	err := serializer.DeserializePath(awsClusterTemplatePath, &spec)
	if err != nil {
		t.Fatal(err)
	}

	amiID, err := spec.resolveAMI()
	if err != nil {
		t.Error(err)
	}

	regex := regexp.MustCompile(`^ami-[a-z0-9]{17}$`)
	if !regex.MatchString(amiID) {
		t.Error("failed to resolve AMI ID.")
	}

	spec.Image = []imageFilter{
		{
			Name:   "name",
			Values: []string{"ami-test"},
		},
	}

	amiID, err = spec.resolveAMI()
	if err == nil {
		t.Error("expected failure to resolve AMI, " +
			"as the filters should resolve more than 1")
	}

	spec.Image = []imageFilter{
		{
			Name:   "name",
			Values: []string{"ami-test"},
		},
		{
			Name:   "owner-id",
			Values: []string{"228170507697"},
		},
	}

	amiID, err = spec.resolveAMI()
	if err != nil {
		t.Error(err)
	}

	regex = regexp.MustCompile(`^ami-[a-z0-9]{17}$`)
	if !regex.MatchString(amiID) {
		t.Error("failed to resolve AMI ID.")
	}
}

func TestCreateAndDestroyAwsCluster(t *testing.T) {
	cloud := getAwsClient(t, awsClusterTemplatePath)
	var spec AwsEnvironment

	err := serializer.DeserializePath(awsClusterTemplatePath, &spec)
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

	cloud.DestroyCluster()
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

func TestDestroyAndDestroySingleNodeAwsCluster(t *testing.T) {
	cloud := getAwsClient(t, awsSingleNodeTemplatePath)
	var spec AwsEnvironment

	err := serializer.DeserializePath(awsSingleNodeTemplatePath, &spec)
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

	cloud.DestroyCluster()
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
