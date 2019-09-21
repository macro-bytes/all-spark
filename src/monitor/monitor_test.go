package monitor

import (
	"cloud"
	"testing"
	"time"
	"util/serializer"
)

func TestHandleCheckinError(t *testing.T) {
	idle := `{
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
"state":"KILLED",
"duration":113949
}],
"status":"ALIVE"
}`
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(idle), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}
}

func TestHandleCheckinDone(t *testing.T) {
	idle := `{
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
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(idle), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusDone {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusDone)
		t.Error("-actual: " + status)
	}
}

func TestHandleCheckinIdle(t *testing.T) {
	idle := `{
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
"completedapps":[],
"status":"ALIVE"
}`
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(idle), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusIdle {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusIdle)
		t.Error("-actual: " + status)
	}
}

func TestHandleCheckinRunning(t *testing.T) {
	running := `{
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
"activeapps":[{
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
"completedapps":[],
"status":"ALIVE"
}`
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(running), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusRunning {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusRunning)
		t.Error("-actual: " + status)
	}
}

func TestPendingTimeoutMonitor(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)

	Run(1, 9999, 9999, 5, 9999)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusPending {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusPending)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 9999, 5, 9999)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + status)
	}
}

func TestIdleTimeoutMonitor(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	setStatus(client.ClusterID, SparkClusterStatusAtEpoch{
		Client:           serlializedClient,
		Timestamp:        time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusIdle,
	})

	Run(1, 9999, 5, 9999, 9999)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusIdle {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusIdle)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 5, 9999, 9999)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + status)
	}
}

func TestMaxRuntime(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	setStatus(client.ClusterID, SparkClusterStatusAtEpoch{
		Client:           serlializedClient,
		Timestamp:        time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusRunning,
	})

	Run(1, 5, 9999, 9999, 9999)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusRunning {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusRunning)
		t.Error("-actual: " + status)
	}

	Run(1, 5, 9999, 9999, 9999)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + status)
	}
}

func TestDoneReportTime(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	setStatus(client.ClusterID, SparkClusterStatusAtEpoch{
		Client:           serlializedClient,
		Timestamp:        time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusDone,
	})

	Run(1, 9999, 9999, 9999, 5)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusDone {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusDone)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 9999, 9999, 5)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + status)
	}
}
