package serializer

import (
	"spark_monitor"
	"strings"
	"template"
	"testing"
)

func TestDeserializePath(t *testing.T) {
	var template template.DockerTemplate

	err := DeserializePath("does-not-exist", &template)
	if err == nil {
		t.Error("Expected non-nil error")
	}

	err = DeserializePath("../../../sample_templates/docker.json", &template)
	if err != nil {
		t.Error(err)
	}
}

func TestSerializeAndDeserialize(t *testing.T) {
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

	var clusterStatus spark_monitor.SparkClusterStatus
	err := Deserialize([]byte(buffer), &clusterStatus)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Serialize(clusterStatus)
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
