package monitor

import (
	"allspark/cloud"
	"allspark/util/serializer"
	"bytes"
	"strconv"
	"testing"
	"time"
)

const (
	IdleStateCheckIn = `{
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

	RunningStateCheckIn = `{
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

	ErrorStateCheckIn = `{
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

	DoneStateCheckIn = `{
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
)

func TestRegisterCluster(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)

	lastKnownStatus, err := getLastEpoch(client.ClusterID)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(lastKnownStatus.Client, serlializedClient) != 0 {
		t.Error("AWS client failed to serialize")
	}

	if lastKnownStatus.CloudEnvironment != cloud.Aws {
		t.Error("cloud environment mismatch")
		t.Error("-expected: " + cloud.Aws)
		t.Error("-actual: " + lastKnownStatus.CloudEnvironment)
	}

	if lastKnownStatus.LastCheckIn == 0 {
		t.Error("last check-in mismatch")
		t.Error("-expected: value > 0")
		t.Error("-actual: " + strconv.FormatInt(lastKnownStatus.LastCheckIn, 10))
	}

	if lastKnownStatus.Timestamp == 0 {
		t.Error("timestamp mismatch")
		t.Error("-expected: value > 0")
		t.Error("-actual: " + strconv.FormatInt(lastKnownStatus.Timestamp, 10))
	}

	if lastKnownStatus.Status != StatusPending {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusPending)
		t.Error("-actual: " + lastKnownStatus.Status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestDuplicateClusterIDHandler(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	DeregisterCluster(client.ClusterID)
	err = RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	if err != nil {
		t.Error(err)
	}

	err = RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	if err == nil {
		t.Error("expected dupicate cluster error")
	}

	DeregisterCluster(client.ClusterID)
}
func TestHandleCheckinError(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(ErrorStateCheckIn), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, "", clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestHandleCheckinDone(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(DoneStateCheckIn), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, "", clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusDone {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusDone)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestHandleCheckinIdle(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(IdleStateCheckIn), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, "", clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusIdle {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusIdle)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestHandleCheckinRunning(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(RunningStateCheckIn), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, "", clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusRunning {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusRunning)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestHandleCheckinAppExitStatus(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	var clusterStatus cloud.SparkClusterStatus
	err = serializer.Deserialize([]byte(RunningStateCheckIn), &clusterStatus)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)
	HandleCheckIn(client.ClusterID, StatusError, clusterStatus)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)

	HandleCheckIn(client.ClusterID, StatusError, clusterStatus)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}

	HandleCheckIn(client.ClusterID, StatusDone, clusterStatus)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestPendingTimeoutMonitor(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)

	Run(1, 9999, 9999, 9999, 5, 9999)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusPending {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusPending)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 9999, 9999, 5, 9999)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}

	DeregisterCluster(client.ClusterID)
}

func TestIdleTimeoutMonitor(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
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
		LastCheckIn:      time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusIdle,
	}, true)

	Run(1, 9999, 5, 9999, 9999, 5)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusIdle {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusIdle)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 5, 9999, 9999, 5)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusDone {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusDone)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 5, 9999, 9999, 5)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + status)
	}
}

func TestMaxRuntime(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
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
		LastCheckIn:      time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusRunning,
	}, true)

	Run(1, 5, 9999, 9999, 9999, 9999)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusRunning {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusRunning)
		t.Error("-actual: " + status)
	}

	Run(1, 5, 9999, 9999, 9999, 9999)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}
}

func TestMaxTimeWithoutCheckin(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
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
		LastCheckIn:      time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusRunning,
	}, true)

	Run(1, 9999, 9999, 5, 9999, 9999)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusRunning {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusRunning)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 9999, 5, 9999, 9999)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusError {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusError)
		t.Error("-actual: " + status)
	}
}

func TestDoneReportTime(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../dist/sample_templates/aws.json", &client)
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
		LastCheckIn:      time.Now().Unix(),
		CloudEnvironment: cloud.Aws,
		Status:           StatusDone,
	}, true)

	Run(1, 9999, 9999, 9999, 9999, 5)
	status := GetLastKnownStatus(client.ClusterID)
	if status != StatusDone {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusDone)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 9999, 9999, 9999, 5)
	status = GetLastKnownStatus(client.ClusterID)
	if status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + status)
	}
}

func TestUnregisterdEpochHandler(t *testing.T) {
	priorStatus, err := getLastEpoch("does-not-exit")
	if err == nil {
		t.Error("expected non-nil error")
	}
	if priorStatus.Status != StatusNotRegistered {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusNotRegistered)
		t.Error("-actual: " + priorStatus.Status)
	}
}
