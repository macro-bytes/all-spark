package monitor

import (
	"cloud"
	"datastore"
	"logger"
	"os"
	"time"
	"util/serializer"
)

// Spark cluster status constants
const (
	StatusNotRegistered = "NOT_REGISTERED"
	StatusPending       = "PENDING"
	StatusIdle          = "IDLE"
	StatusRunning       = "RUNNING"
	StatusDone          = "DONE"
	StatusError         = "ERROR"
	StatusFinished      = "FINISHED"
	statusMap           = "STATUS_MAP"
	monitorLock         = "MONITOR_LOCK"
)

// SparkClusterStatusAtEpoch describes the state of a cluster
// at a given timestamp
type SparkClusterStatusAtEpoch struct {
	Timestamp        int64
	Status           string
	Client           []byte
	CloudEnvironment string
}

// GetClientData - Returns the serialized and cloud environment
func GetClientData(clusterID string) ([]byte, string, error) {
	state, err := getLastEpoch(clusterID)
	if err != nil {
		logger.GetError().Printf("Unable to retrieve state for cluster %v", clusterID)
		return nil, "", err
	}

	return state.Client, state.CloudEnvironment, nil
}

// HandleCheckIn - handles spark monitor check-in http requests
func HandleCheckIn(clusterID string, clusterStatus cloud.SparkClusterStatus) {
	logger.GetInfo().Printf("cluster: %v, status: %+v", clusterID, clusterStatus)

	priorClusterState, err := getLastEpoch(clusterID)
	if err != nil {
		logger.GetError().Println(err)
	}

	epochStatus := SparkClusterStatusAtEpoch{
		Timestamp:        getTimestamp(),
		Status:           getReportedStatus(clusterStatus),
		Client:           priorClusterState.Client,
		CloudEnvironment: priorClusterState.CloudEnvironment,
	}

	setStatus(clusterID, epochStatus)
}

// RegisterCluster - registers newly created spark
// cluster with a pending status
func RegisterCluster(clusterID string, cloudEnvironment string, serializedClient []byte) {
	logger.GetInfo().Printf("registering cluster: %s, %s, %s",
		clusterID, cloudEnvironment, serializedClient)
	setStatus(clusterID, SparkClusterStatusAtEpoch{
		Status:           StatusPending,
		Timestamp:        getTimestamp(),
		Client:           serializedClient,
		CloudEnvironment: cloudEnvironment,
	})
}

// DeregisterCluster - registers newly created spark
// cluster with a pending status
func DeregisterCluster(clusterID string) {
	logger.GetInfo().Printf("deregistering cluster %s", clusterID)
	client := datastore.GetRedisClient()
	defer client.Close()

	client.HDel(statusMap, clusterID)
}

func getReportedStatus(status cloud.SparkClusterStatus) string {
	if len(status.ActiveApps) > 0 {
		return StatusRunning
	} else if len(status.CompletedApps) > 0 {
		if status.CompletedApps[0].State != StatusFinished {
			return StatusError
		}
		return StatusDone
	}

	return StatusIdle
}

func getLastEpoch(clusterID string) (SparkClusterStatusAtEpoch, error) {
	client := datastore.GetRedisClient()
	defer client.Close()

	var clusterState SparkClusterStatusAtEpoch
	err := serializer.Deserialize([]byte(client.HGet(statusMap, clusterID).Val()),
		&clusterState)
	if err != nil {
		return clusterState, err
	}

	return clusterState, nil
}

// GetLastKnownStatus - returns the last known status of the cluster
func GetLastKnownStatus(clusterID string) string {
	clusterState, err := getLastEpoch(clusterID)
	if err != nil {
		return StatusNotRegistered
	}
	return clusterState.Status
}

func setStatus(clusterID string, status SparkClusterStatusAtEpoch) {
	logger.GetInfo().Printf("setting status %s, status: %+v", clusterID, status.Status)
	client := datastore.GetRedisClient()
	defer client.Close()

	result, err := serializer.Serialize(status)
	if err != nil {
		logger.GetError().Println(err)
	}

	client.HSet(statusMap, clusterID, string(result))
}

// Run - daemon used for monitoring all spark clusters;
// monitor will run for the specified number of iterations, or indefinitely
// if iterations <= 0.
func Run(iterations int, maxRuntime int64, idleTimeout int64,
	pendingTimeout int64, doneReportTime int64) {

	if iterations <= 0 {
		for {
			if acquireLock() {
				logger.GetInfo().Println("acquired lock")
				monitorClusterHelper(maxRuntime, idleTimeout,
					pendingTimeout, doneReportTime)
				releaseLock()
			}
			time.Sleep(10 * time.Second)
		}
	}
	for i := 0; i < iterations; i++ {
		if acquireLock() {
			monitorClusterHelper(maxRuntime, idleTimeout,
				pendingTimeout, doneReportTime)
			releaseLock()
		}
		time.Sleep(10 * time.Second)
	}
}

func monitorClusterHelper(maxRuntime int64, idleTimeout int64,
	pendingTimeout int64, doneReportTime int64) {

	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()

	for clusterID, buffer := range redisClient.HGetAll(statusMap).Val() {
		var status SparkClusterStatusAtEpoch
		serializer.Deserialize([]byte(buffer), &status)

		client, err := cloud.Create(status.CloudEnvironment, status.Client)
		if err != nil {
			logger.GetError().Println(err)
		}

		currentTime := getTimestamp()
		switch status.Status {
		case StatusPending:
			logger.GetInfo().Printf("monitor reported %s for cluster %s",
				StatusPending, clusterID)
			if currentTime-status.Timestamp > pendingTimeout {
				client.DestroyCluster()
				DeregisterCluster(clusterID)
			}
			break
		case StatusIdle:
			logger.GetInfo().Printf("monitor reported %s for cluster %s",
				StatusIdle, clusterID)
			if currentTime-status.Timestamp > idleTimeout {
				client.DestroyCluster()
				DeregisterCluster(clusterID)
			}
			break
		case StatusRunning:
			logger.GetInfo().Printf("monitor reported %s for cluster %s",
				StatusRunning, clusterID)
			if currentTime-status.Timestamp > maxRuntime {
				client.DestroyCluster()
				DeregisterCluster(clusterID)
			}
			break
		case StatusDone:
			logger.GetInfo().Printf("monitor reported %s for cluster %s",
				StatusDone, clusterID)
			client.DestroyCluster()
			if currentTime-status.Timestamp > doneReportTime {
				DeregisterCluster(clusterID)
			}
			break
		default:
			logger.GetInfo().Printf("monitor reported no status for cluster %s",
				clusterID)
			break
		}
	}
}

func releaseLock() {
	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()

	redisClient.Del(monitorLock)
}

func acquireLock() bool {
	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()

	id, err := os.Hostname()
	if err != nil {
		logger.GetError().Println(err)
		return false
	}

	return redisClient.SetNX(monitorLock, id, 15*time.Minute).Val()
}

func getTimestamp() int64 {
	return time.Now().Unix()
}
