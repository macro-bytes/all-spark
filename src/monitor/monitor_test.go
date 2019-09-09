package monitor

import (
	"cloud"
	"testing"
	"time"
	"util/serializer"
)

func TestPendingTimeoutMonitor(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../sample_templates/aws.json", &client)
	if err != nil {
		t.Error(err)
	}

	serlializedClient, err := serializer.Serialize(client)
	if err != nil {
		t.Error(err)
	}

	RegisterCluster(client.ClusterID, cloud.Aws, serlializedClient)

	Run(1, 9999, 9999, 5)
	status := getPriorStatus(client.ClusterID)
	if status != StatusPending {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusPending)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 9999, 5)
	status = getPriorStatus(client.ClusterID)
	if status != StatusUnknown {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusUnknown)
		t.Error("-actual: " + status)
	}
}

func TestIdleTimeoutMonitor(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../sample_templates/aws.json", &client)
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

	Run(1, 9999, 5, 9999)
	status := getPriorStatus(client.ClusterID)
	if status != StatusIdle {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusIdle)
		t.Error("-actual: " + status)
	}

	Run(1, 9999, 5, 9999)
	status = getPriorStatus(client.ClusterID)
	if status != StatusUnknown {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusUnknown)
		t.Error("-actual: " + status)
	}
}

func TestMaxRuntime(t *testing.T) {
	var client cloud.AwsEnvironment
	err := serializer.DeserializePath("../../sample_templates/aws.json", &client)
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

	Run(1, 5, 9999, 9999)
	status := getPriorStatus(client.ClusterID)
	if status != StatusRunning {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusRunning)
		t.Error("-actual: " + status)
	}

	Run(1, 5, 9999, 9999)
	status = getPriorStatus(client.ClusterID)
	if status != StatusUnknown {
		t.Error("status mismatch")
		t.Error("-expected: " + StatusUnknown)
		t.Error("-actual: " + status)
	}
}
