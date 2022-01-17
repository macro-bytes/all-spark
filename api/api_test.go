package api

import (
	"allspark/cloud"
	"allspark/monitor"
	"allspark/util/serializer"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	awsTemplatePath    = "../dist/sample_templates/aws.json"
	dockerTemplatePath = "../dist/sample_templates/docker.json"
	azureTemplatePath  = "../dist/sample_templates/azure.json"
)

func getAwsClient(t *testing.T) cloud.CloudEnvironment {
	templateConfig, err := cloud.ReadTemplateConfiguration(awsTemplatePath)
	if err != nil {
		t.Fatal(err)
	}

	client, err := cloud.Create(cloud.Aws, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func getDockerClient(t *testing.T) cloud.CloudEnvironment {
	templateConfig, err := cloud.ReadTemplateConfiguration(dockerTemplatePath)
	if err != nil {
		t.Fatal(err)
	}

	client, err := cloud.Create(cloud.Docker, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func getAzureClient(t *testing.T) cloud.CloudEnvironment {
	templateConfig, err := cloud.ReadTemplateConfiguration(azureTemplatePath)
	if err != nil {
		t.Fatal(err)
	}

	client, err := cloud.Create(cloud.Azure, templateConfig)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func getBadCreateFormDataAws() []byte {
	var template = cloud.AwsEnvironment{
		ClusterID:     "test",
		EBSVolumeSize: 0,
		IAMRole:       "test",
	}

	buff, _ := json.Marshal(template)
	return buff
}

func getValidCreateFormDataAws() []byte {
	var template cloud.AwsEnvironment
	serializer.DeserializePath(awsTemplatePath,
		&template)

	buff, _ := json.Marshal(template)
	return buff
}

func getBadCreateFormDataDocker() []byte {
	var template = cloud.DockerEnvironment{
		ClusterID: "test",
		Image:     "image-does-not-exist",
	}

	buff, _ := json.Marshal(template)
	return buff
}

func getValidCreateFormDataDocker() []byte {
	var template cloud.DockerEnvironment
	serializer.DeserializePath(dockerTemplatePath,
		&template)

	buff, _ := json.Marshal(template)
	return buff
}

func getBadCreateFormDataAzure() []byte {
	var template = cloud.AzureEnvironment{
		ClusterID: "test",
		ImageBlob: "does-not-exist",
	}

	buff, _ := json.Marshal(template)
	return buff
}

func getValidCreateFormDataAzure() []byte {
	var template cloud.AzureEnvironment
	serializer.DeserializePath(azureTemplatePath,
		&template)

	buff, _ := json.Marshal(template)
	return buff
}

func getDestroyClusterFormDocker() string {
	var template cloud.DockerEnvironment
	serializer.DeserializePath(dockerTemplatePath,
		&template)
	formData := url.Values{}
	formData.Set("clusterID", template.ClusterID)
	return formData.Encode()
}

func getDestroyClusterFormAws() string {
	var template cloud.AwsEnvironment
	serializer.DeserializePath(awsTemplatePath,
		&template)
	formData := url.Values{}
	formData.Set("clusterID", template.ClusterID)
	return formData.Encode()
}

func getDestroyClusterFormAzure() string {
	var template cloud.AzureEnvironment
	serializer.DeserializePath(azureTemplatePath,
		&template)
	formData := url.Values{}
	formData.Set("clusterID", template.ClusterID)
	return formData.Encode()
}

func testHTTPRequest(t *testing.T,
	handlerFunction func(http.ResponseWriter,
		*http.Request), method string,
	route string, body io.Reader, expectedStatusCode int,
	formURLEncode bool) {

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunction)

	req, err := http.NewRequest(method, route, body)
	if formURLEncode {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != expectedStatusCode {
		t.Fatalf("unexpected status code: got %v, expected %v",
			status, expectedStatusCode)
	}
}

func TestHealthCheck(t *testing.T) {
	testHTTPRequest(t, healthCheck, "POST", "/healthCheck",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, healthCheck, "GET", "/healthCheck",
		nil, http.StatusOK, false)
}

func TestCheckin(t *testing.T) {
	testHTTPRequest(t, checkIn, "GET", "/check-in",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, checkIn, "POST", "/check-in",
		nil, http.StatusBadRequest, false)

	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}
	client.ClusterID = "local"

	monitor.RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)

	idle_cluster_state := []byte(`{
		"ClusterID": "local",
		"AppExitStatus": "",
		"Status": {	"url" : "spark://simulated-local-mode-cluster:7077",
			"workers" : [ ],
			"aliveworkers" : 0,
			"cores" : 0,
			"coresused" : 0,
			"memory" : 0,
			"memoryused" : 0,
			"activeapps" : [ ],
			"completedapps" : [ ],
			"activedrivers" : [ ],
			"completeddrivers" : [ ],
			"status" : "ALIVE"
		}
	}`)

	testHTTPRequest(t, checkIn, "POST", "/check-in",
		bytes.NewReader(idle_cluster_state), http.StatusOK, false)
	clusterStatus := monitor.GetLastKnownStatus("local")
	if clusterStatus != monitor.StatusIdle {
		t.Error("Expected cluster status " + monitor.StatusIdle + ", got " + clusterStatus)
	}

	running_cluster_state := []byte(`{
		"ClusterID": "local",
		"AppExitStatus": "",
		"Status": {
			"url" : "spark://simulated-local-mode-cluster::7077",
			"workers" : [ ],
			"aliveworkers" : 0,
			"cores" : 0,
			"coresused" : 0,
			"memory" : 0,
			"memoryused" : 0,
			"activeapps" : [ {
				"id" : "app-0-0000",
				"starttime" : 0,
				"name" : "pyspark-shell",
				"cores" : 0,
				"user" : "spark-user",
				"memoryperslave" : 0,
				"submitdate" : "Thu Jan 13 17:47:31 GMT 2022",
				"state" : "WAITING",
				"duration" : 0
			} ],
			"completedapps" : [ ],
			"activedrivers" : [ ],
			"completeddrivers" : [ ],
			"status" : "ALIVE"
		}
	}`)

	testHTTPRequest(t, checkIn, "POST", "/check-in",
		bytes.NewReader(running_cluster_state), http.StatusOK, false)
	clusterStatus = monitor.GetLastKnownStatus("local")
	if clusterStatus != monitor.StatusRunning {
		t.Error("Expected cluster status " + monitor.StatusRunning + ", got " + clusterStatus)
	}

	testHTTPRequest(t, checkIn, "POST", "/check-in",
		bytes.NewReader(idle_cluster_state), http.StatusOK, false)
	clusterStatus = monitor.GetLastKnownStatus("local")
	if clusterStatus != monitor.StatusIdle {
		t.Error("Expected cluster status " + monitor.StatusIdle + ", got " + clusterStatus)
	}

	time.Sleep(2 * time.Second)
	monitor.Run(1, 9999, 1, 9999, 9999, 9999, 9999)

	testHTTPRequest(t, checkIn, "POST", "/check-in",
		bytes.NewReader(idle_cluster_state), http.StatusOK, false)
	clusterStatus = monitor.GetLastKnownStatus("local")
	if clusterStatus != monitor.StatusDone {
		t.Error("Expected cluster status " + monitor.StatusDone + ", got " + clusterStatus)
	}

	monitor.DeregisterCluster(client.ClusterID)

}

func TestGetStatus(t *testing.T) {
	testHTTPRequest(t, getStatus, "POST", "/getStatus",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, getStatus, "GET", "/getStatus",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, getStatus, "GET", "/getStatus?clusterID=test",
		nil, http.StatusOK, false)
}

func TestCreateAndDestroyClusterAWS(t *testing.T) {
	testHTTPRequest(t, createClusterAws, "GET", "/aws/create",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterAws, "POST", "/aws/create",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterAws, "POST", "/aws/create",
		bytes.NewReader(getBadCreateFormDataAws()), http.StatusBadRequest,
		false)
	testHTTPRequest(t, createClusterAws, "POST", "/aws/create",
		bytes.NewReader(getValidCreateFormDataAws()), http.StatusOK, false)

	testHTTPRequest(t, terminateDocker, "POST", "/docker/terminate",
		strings.NewReader(getDestroyClusterFormAws()),
		http.StatusBadRequest, true)
	testHTTPRequest(t, terminateDocker, "POST", "/azure/terminate",
		strings.NewReader(getDestroyClusterFormAws()),
		http.StatusBadRequest, true)
	testHTTPRequest(t, terminateAws, "POST", "/aws/terminate",
		strings.NewReader(getDestroyClusterFormAws()),
		http.StatusOK, true)
	testHTTPRequest(t, terminateAws, "GET", "/aws/terminate",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, terminateAws, "POST", "/aws/terminate",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, terminateAws, "POST",
		"/aws/terminate", bytes.NewReader(getBadCreateFormDataAws()),
		http.StatusBadRequest, false)

	getAwsClient(t).DestroyCluster()
}

func TestCreateAndDestroyClusterAzure(t *testing.T) {
	testHTTPRequest(t, createClusterAzure, "GET", "/azure/create",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterAzure, "POST", "/azure/create",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterAzure, "POST", "/azure/create",
		bytes.NewReader(getBadCreateFormDataAzure()), http.StatusBadRequest,
		false)
	testHTTPRequest(t, createClusterAzure, "POST", "/azure/create",
		bytes.NewReader(getValidCreateFormDataAzure()), http.StatusOK, false)

	testHTTPRequest(t, terminateDocker, "POST", "/docker/terminate",
		strings.NewReader(getDestroyClusterFormAzure()),
		http.StatusBadRequest, true)
	testHTTPRequest(t, terminateAws, "POST", "/aws/terminate",
		strings.NewReader(getDestroyClusterFormAzure()),
		http.StatusBadRequest, true)
	testHTTPRequest(t, terminateAzure, "POST", "/azure/terminate",
		strings.NewReader(getDestroyClusterFormAzure()),
		http.StatusOK, true)
	testHTTPRequest(t, terminateAzure, "GET", "/azure/terminate",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, terminateAzure, "POST", "/azure/terminate",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, terminateAzure, "POST",
		"/azure/terminate", bytes.NewReader(getBadCreateFormDataAzure()),
		http.StatusBadRequest, false)

	getAzureClient(t).DestroyCluster()
	time.Sleep(5 * time.Minute)
}

func TestCreateAndDestroyClusterDocker(t *testing.T) {
	testHTTPRequest(t, createClusterDocker, "GET",
		"/docker/create", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterDocker, "POST",
		"/docker/create", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterDocker, "POST",
		"/docker/create", bytes.NewReader(getBadCreateFormDataDocker()),
		http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterDocker, "POST",
		"/docker/create",
		bytes.NewReader(getValidCreateFormDataDocker()), http.StatusOK, false)

	testHTTPRequest(t, terminateAws, "POST", "/aws/terminate",
		strings.NewReader(getDestroyClusterFormDocker()),
		http.StatusBadRequest, true)
	testHTTPRequest(t, terminateAws, "POST", "/azure/terminate",
		strings.NewReader(getDestroyClusterFormDocker()),
		http.StatusBadRequest, true)
	testHTTPRequest(t, terminateDocker, "POST", "/docker/terminate",
		strings.NewReader(getDestroyClusterFormDocker()),
		http.StatusOK, true)
	testHTTPRequest(t, terminateDocker, "GET",
		"/docker/terminate", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, terminateDocker, "POST",
		"/docker/terminate", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, terminateDocker, "POST",
		"/docker/terminate", bytes.NewReader(getBadCreateFormDataDocker()),
		http.StatusBadRequest, false)

	getDockerClient(t).DestroyCluster()
}
