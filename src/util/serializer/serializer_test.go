package serializer_test

import (
	"cloud"
	"monitor"
	"reflect"
	"strings"
	"testing"
	"time"
	"util/serializer"
)

func TestDeserializePath(t *testing.T) {
	var template cloud.DockerEnvironment
	err := serializer.DeserializePath("does-not-exist", &template)
	if err == nil {
		t.Error("Expected non-nil error")
	}

	err = serializer.DeserializePath("../../../dist/sample_templates/docker.json", &template)
	if err != nil {
		t.Error(err)
	}
}

func TestSerializeAndDeserializeEpochStatus(t *testing.T) {
	var dummyClient cloud.AwsEnvironment
	err := serializer.DeserializePath("../../../dist/sample_templates/aws.json", &dummyClient)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(dummyClient)
	if err != nil {
		t.Error(err)
	}

	state := monitor.SparkClusterStatusAtEpoch{
		Client:           serlializedClient,
		Status:           monitor.StatusPending,
		CloudEnvironment: cloud.Aws,
		Timestamp:        time.Now().Unix(),
	}

	serializedState, err := serializer.Serialize(state)
	if err != nil {
		t.Error(err)
	}

	var stateTest monitor.SparkClusterStatusAtEpoch
	serializer.Deserialize(serializedState, &stateTest)

	if !reflect.DeepEqual(stateTest, state) {
		t.Error("deserialization failed")
	}

	var clientTest cloud.AwsEnvironment
	err = serializer.Deserialize(stateTest.Client, &clientTest)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(clientTest, dummyClient) {
		t.Error("deserialization failed")
	}
}

func TestSerializeAndDeserializeClusterStatus(t *testing.T) {
	buffer := `{
"url":"spark://ip-172-30-0-100.us-west-2.compute.internal:7077",
"workers":[{
"id":"worker-20190904193157-172.30.0.132-7078",
"host":"172.30.0.132",
"port":7078,
"webuiaddress":"http://172.30.0.132:8081",
"cores":8,
"coresused":0,
"coresfree":8,
"memory":30348,
"memoryused":0,
"memoryfree":30348,
"state":"ALIVE",
"lastheartbeat":1567625653486
},{
"id":"worker-20190904193157-172.30.0.12-7078",
"host":"172.30.0.12",
"port":7078,
"webuiaddress":"http://172.30.0.12:8081",
"cores":8,
"coresused":0,
"coresfree":8,
"memory":30348,
"memoryused":0,
"memoryfree":30348,
"state":"ALIVE",
"lastheartbeat":1567625653258
}],
"aliveworkers":2,
"cores":16,
"coresused":0,
"memory":60696,
"memoryused":0,
"activeapps":[],
"completedapps":[{
"id":"app-20190904193210-0000",
"starttime":1567625530255,
"name":"Sparkling Water Driver",
"cores":16,
"user":"root",
"memoryperslave":1024,
"submitdate":"Wed Sep 04 19:32:10 GMT 2019",
"state":"FINISHED",
"duration":113949
}],
"status":"ALIVE"
}`

	var clusterStatus cloud.SparkClusterStatus
	err := serializer.Deserialize([]byte(buffer), &clusterStatus)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := serializer.Serialize(clusterStatus)
	if err != nil {
		t.Fatal(err)
	}

	expected := strings.ReplaceAll(buffer, "\n", "")
	if string(actual) != expected {
		t.Error("serialization/deserialization failed")
		t.Error("-expected: " + expected)
		t.Fatal("-actual: " + string(actual))
	}
}
