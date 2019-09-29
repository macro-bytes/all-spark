package api

import (
	"bytes"
	"cloud"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"util/serializer"
)

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
	serializer.DeserializePath("../../dist/sample_templates/aws.json",
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
	serializer.DeserializePath("../../dist/sample_templates/docker.json",
		&template)

	buff, _ := json.Marshal(template)
	return buff
}

func getDestroyClusterFormDocker() string {
	var template cloud.DockerEnvironment
	serializer.DeserializePath("../../dist/sample_templates/docker.json",
		&template)
	formData := url.Values{}
	formData.Set("clusterID", template.ClusterID)
	return formData.Encode()
}

func getDestroyClusterFormAws() string {
	var template cloud.AwsEnvironment
	serializer.DeserializePath("../../dist/sample_templates/aws.json",
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
		t.Errorf("unexpected status code: got %v, expected %v",
			status, expectedStatusCode)
	}
}

func TestCheckin(t *testing.T) {
	testHTTPRequest(t, checkIn, "GET", "/checkIn",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, checkIn, "POST", "/checkIn",
		nil, http.StatusBadRequest, false)
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
	testHTTPRequest(t, createClusterAws, "GET", "/aws/createCluster",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterAws, "POST", "/aws/createCluster",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterAws, "POST", "/aws/createCluster",
		bytes.NewReader(getBadCreateFormDataAws()), http.StatusBadRequest,
		false)
	testHTTPRequest(t, createClusterAws, "POST", "/aws/createCluster",
		bytes.NewReader(getValidCreateFormDataAws()), http.StatusOK, false)

	testHTTPRequest(t, destroyCluster, "GET", "/destroyCluster",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, destroyCluster, "POST", "/destroyCluster",
		nil, http.StatusBadRequest, false)
	testHTTPRequest(t, destroyCluster, "POST",
		"/destroyCluster", bytes.NewReader(getBadCreateFormDataAws()),
		http.StatusBadRequest, false)

	testHTTPRequest(t, destroyCluster, "POST", "/destroyCluster",
		strings.NewReader(getDestroyClusterFormAws()),
		http.StatusOK, true)
}

func TestCreateAndDestroyClusterDocker(t *testing.T) {
	testHTTPRequest(t, createClusterDocker, "GET",
		"/docker/createCluster", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterDocker, "POST",
		"/docker/createCluster", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterDocker, "POST",
		"/docker/createCluster", bytes.NewReader(getBadCreateFormDataDocker()),
		http.StatusBadRequest, false)
	testHTTPRequest(t, createClusterDocker, "POST",
		"/docker/createCluster",
		bytes.NewReader(getValidCreateFormDataDocker()), http.StatusOK, false)

	testHTTPRequest(t, destroyCluster, "GET",
		"/destroyCluster", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, destroyCluster, "POST",
		"/destroyCluster", nil, http.StatusBadRequest, false)
	testHTTPRequest(t, destroyCluster, "POST",
		"/destroyCluster", bytes.NewReader(getBadCreateFormDataDocker()),
		http.StatusBadRequest, false)

	testHTTPRequest(t, destroyCluster, "POST", "/destroyCluster",
		strings.NewReader(getDestroyClusterFormDocker()), http.StatusOK, true)
}
