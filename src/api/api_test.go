package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"template"
	"testing"
	"util/template_reader"
)

func createBadFormDataAws() []byte {
	var template = template.AwsTemplate{
		ClusterID:     "test",
		EBSVolumeSize: 0,
		IAMRole:       "test",
	}

	buff, _ := json.Marshal(template)
	return buff
}

func createValidFormDataAws() []byte {
	var template template.AwsTemplate
	template_reader.Deserialize("../../sample_templates/aws.json",
		&template)

	buff, _ := json.Marshal(template)
	return buff
}

func createBadFormDataDocker() []byte {
	var template = template.DockerTemplate{
		ClusterID: "test",
		Image:     "image-does-not-exist",
	}

	buff, _ := json.Marshal(template)
	return buff
}

func createValidFormDataDocker() []byte {
	var template template.AwsTemplate
	template_reader.Deserialize("../../sample_templates/docker.json",
		&template)

	buff, _ := json.Marshal(template)
	return buff
}

func testHTTPRequest(t *testing.T,
	handlerFunction func(http.ResponseWriter,
		*http.Request), method string,
	route string, body io.Reader, expectedStatusCode int) {

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunction)

	req, err := http.NewRequest(method, route, body)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("unexpected status code: got %v, expected %v",
			status, expectedStatusCode)
	}
}

func TestCreateAndDestroyClusterAWS(t *testing.T) {
	testHTTPRequest(t, createClusterAws, "GET",
		"/aws/createCluster", nil, http.StatusBadRequest)
	testHTTPRequest(t, createClusterAws, "POST",
		"/aws/createCluster", nil, http.StatusBadRequest)
	testHTTPRequest(t, createClusterAws, "POST",
		"/aws/createCluster",
		bytes.NewReader(createBadFormDataAws()), http.StatusBadRequest)
	testHTTPRequest(t, createClusterAws, "POST",
		"/aws/createCluster",
		bytes.NewReader(createValidFormDataAws()), http.StatusOK)

	testHTTPRequest(t, destroyClusterAws, "GET",
		"/aws/destroyCluster", nil, http.StatusBadRequest)
	testHTTPRequest(t, destroyClusterAws, "POST",
		"/aws/destroyCluster", nil, http.StatusBadRequest)
	testHTTPRequest(t, destroyClusterAws, "POST",
		"/aws/destroyCluster",
		bytes.NewReader(createBadFormDataAws()), http.StatusBadRequest)
	testHTTPRequest(t, destroyClusterAws, "POST",
		"/aws/destroyCluster",
		bytes.NewReader(createValidFormDataAws()), http.StatusOK)
}
