package monitor

import (
	"allspark/cloud"
	"allspark/datastore"
	"allspark/logger"
	"errors"
	"os"
	"strconv"
	"time"

	"allspark/util/serializer"
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
	StatusCanceled      = "CANCELED"
	StatusTerminating   = "TERMINATING"
	statusMap           = "STATUS_MAP"
	monitorLock         = "MONITOR_LOCK"
	clusterLockPreifx   = "cluster.lock."
)

// SparkClusterStatusAtEpoch describes the state of a cluster
// at a given timestamp
type SparkClusterStatusAtEpoch struct {
	LastCheckIn      int64
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
func HandleCheckIn(clusterID string, appExitStatus string,
	clusterStatus cloud.SparkClusterStatus) {

	err := acquireClusterLock(clusterID, "check-in", 5)
	if err != nil {
		logger.GetError().Println(err)
		return
	}
	defer releaseClusterLock(clusterID)

	logger.GetInfo().Printf("cluster: %v, app exit status: %v, status: %+v",
		clusterID, appExitStatus, clusterStatus)

	priorClusterState, err := getLastEpoch(clusterID)
	if err != nil {
		logger.GetInfo().Printf("cluster: %v appears to have been terminated and deregistered",
			clusterID)
	}

	logger.GetInfo().Printf("cluster: %v prior cluster state: %v",
		clusterID, priorClusterState.Status)

	if priorClusterState.Status == StatusNotRegistered {
		logger.GetInfo().Printf("cluster: %v checked-in, but is not currently"+
			" registered and likely set for termination",
			clusterID)
		return
	}

	if priorClusterState.Status == StatusCanceled {
		logger.GetInfo().Printf("cluster: %v checked-in, but is currently"+
			" set for cancelation",
			clusterID)
		return
	}

	var timestamp int64
	reportedStatus := getReportedStatus(appExitStatus, clusterStatus)
	if reportedStatus == StatusError {
		logger.GetError().Printf("cluster: %v reported status: %+v", clusterID, StatusError)
	} else {
		logger.GetInfo().Printf("cluster: %v reported status: %+v", clusterID, reportedStatus)
	}

	if reportedStatus != priorClusterState.Status {
		timestamp = getTimestamp()
	} else {
		timestamp = priorClusterState.Timestamp
	}

	epochStatus := SparkClusterStatusAtEpoch{
		LastCheckIn:      getTimestamp(),
		Timestamp:        timestamp,
		Status:           reportedStatus,
		Client:           priorClusterState.Client,
		CloudEnvironment: priorClusterState.CloudEnvironment,
	}

	if priorClusterState.Status != StatusDone &&
		priorClusterState.Status != StatusError {
		setStatus(clusterID, epochStatus, true)
	}

	if priorClusterState.Status == StatusDone &&
		reportedStatus == StatusError {
		setStatus(clusterID, epochStatus, true)
	}
}

// RegisterCluster - registers newly created spark
// cluster with a pending status
func RegisterCluster(clusterID string, cloudEnvironment string, serializedClient []byte) error {
	logger.GetInfo().Printf("registering cluster: %s, %s",
		clusterID, cloudEnvironment)

	success := setStatus(clusterID, SparkClusterStatusAtEpoch{
		Status:           StatusPending,
		Timestamp:        getTimestamp(),
		LastCheckIn:      getTimestamp(),
		Client:           serializedClient,
		CloudEnvironment: cloudEnvironment,
	}, false)

	if !success {
		return errors.New("cluster" + clusterID + " already exists")
	}

	return nil
}

// DeregisterCluster - registers newly created spark
// cluster with a pending status
func DeregisterCluster(clusterID string) {
	logger.GetInfo().Printf("deregistering cluster %s", clusterID)
	client := datastore.GetRedisClient()
	defer client.Close()

	client.HDel(statusMap, clusterID)
}

func getReportedStatus(appExitStatus string, status cloud.SparkClusterStatus) string {
	if len(appExitStatus) > 0 {
		// currently all appExitStates with length > 0 are assumed to be error states
		return StatusError
	}

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
		clusterState.Status = StatusNotRegistered
		return clusterState, err
	}

	if clusterState.Status == "" {
		clusterState.Status = StatusNotRegistered
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

// SetCanceled - Sets the cluster to StatusCanceled so be terminated
func SetCanceled(clusterID string) error {
	logger.GetInfo().Printf("handling request to cancel cluster %v ",
		clusterID)
	err := acquireClusterLock(clusterID, "canceled", 5)
	if err != nil {
		logger.GetError().Println(err)
		return err
	}
	defer releaseClusterLock(clusterID)

	priorClusterState, err := getLastEpoch(clusterID)
	if err != nil {
		logger.GetInfo().Printf("cluster: %v appears to have been terminated and deregistered",
			clusterID)
		return err
	}

	if priorClusterState.Status != StatusDone &&
		priorClusterState.Status != StatusError &&
		priorClusterState.Status != StatusCanceled &&
		priorClusterState.Status != StatusNotRegistered {

		epochStatus := SparkClusterStatusAtEpoch{
			LastCheckIn:      priorClusterState.LastCheckIn,
			Timestamp:        getTimestamp(),
			Status:           StatusCanceled,
			Client:           priorClusterState.Client,
			CloudEnvironment: priorClusterState.CloudEnvironment,
		}

		setStatus(clusterID, epochStatus, true)
	} else {
		logger.GetInfo().Printf("cluster %v with status %v will not be set to canceled",
			clusterID, priorClusterState.Status)
	}

	return nil
}

func setStatus(clusterID string, status SparkClusterStatusAtEpoch, overwrite bool) bool {
	logger.GetInfo().Printf("setting status %s, status: %+v", clusterID, status.Status)
	client := datastore.GetRedisClient()
	defer client.Close()

	result, err := serializer.Serialize(status)
	if err != nil {
		logger.GetError().Println(err)
	}

	if overwrite {
		return client.HSet(statusMap, clusterID, string(result)).Val()
	}

	return client.HSetNX(statusMap, clusterID, string(result)).Val()
}

// Run - daemon used for monitoring all spark clusters;
// monitor will run for the specified number of iterations, or indefinitely
// if iterations <= 0.
func Run(iterations int, maxRuntime int64, idleTimeout int64,
	maxTimeWithoutCheckin int64, pendingTimeout int64,
	doneReportTime int64, cancelTerminationDelay int64) {

	if iterations <= 0 {
		for {
			if acquireMonitorLock() {
				logger.GetDebug().Println("acquired lock")
				monitorClusterHelper(maxRuntime, idleTimeout,
					maxTimeWithoutCheckin, pendingTimeout,
					doneReportTime, cancelTerminationDelay)
				releaseLock()
			}
			time.Sleep(10 * time.Second)
		}
	}
	for i := 0; i < iterations; i++ {
		if acquireMonitorLock() {
			monitorClusterHelper(maxRuntime, idleTimeout,
				maxTimeWithoutCheckin, pendingTimeout,
				doneReportTime, cancelTerminationDelay)
			releaseLock()
		}
		time.Sleep(10 * time.Second)
	}
}

func terminateCluster(client cloud.CloudEnvironment) {
	err := client.DestroyCluster()
	if err != nil {
		logger.GetError().Println(err)
	}
}

func monitorClusterHelper(maxRuntime int64, idleTimeout int64,
	maxTimeWithoutCheckin int64, pendingTimeout int64,
	doneReportTime int64, cancelTerminationDelay int64) {

	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()

	for clusterID, buffer := range redisClient.HGetAll(statusMap).Val() {
		err := acquireClusterLock(clusterID, "canceled", 5)
		if err != nil {
			logger.GetError().Println(err)
			continue
		}

		var status SparkClusterStatusAtEpoch
		serializer.Deserialize([]byte(buffer), &status)

		client, err := cloud.Create(status.CloudEnvironment, status.Client)
		if err != nil {
			logger.GetError().Println(err)
			logger.GetError().Printf("cluster does not appear to be valid %v: %v",
				clusterID, redisClient.HGet(statusMap, clusterID).Val())
			logger.GetError().Printf("deregistering cluster %v", clusterID)
			DeregisterCluster(clusterID)
		} else {
			currentTime := getTimestamp()
			if currentTime-status.LastCheckIn > maxTimeWithoutCheckin &&
				status.Status != StatusDone && status.Status != StatusError &&
				status.Status != StatusPending {
				logger.GetError().Printf("max time without check-in exceeded for cluster %s; terminating",
					clusterID)

				status.Status = StatusError
				status.Timestamp = getTimestamp()
				setStatus(clusterID, status, true)
			} else if currentTime-status.Timestamp > maxRuntime {
				logger.GetError().Printf("max run-time exceeded for cluster %s; terminating",
					clusterID)
				status.Status = StatusError
				status.Timestamp = getTimestamp()
				setStatus(clusterID, status, true)
			} else {
				switch status.Status {
				case StatusPending:
					logger.GetInfo().Printf("monitor reported %s for cluster %s",
						status.Status, clusterID)
					if currentTime-status.Timestamp > pendingTimeout {
						logger.GetError().Printf("pending timeout exceeded for cluster %s; terminating",
							clusterID)

						status.Status = StatusError
						status.Timestamp = getTimestamp()
						setStatus(clusterID, status, true)
					}
					break
				case StatusIdle:
					logger.GetInfo().Printf("monitor reported %s for cluster %s",
						status.Status, clusterID)
					if currentTime-status.Timestamp > idleTimeout {
						logger.GetInfo().Printf("idle timeout exceeded for cluster %s; terminating",
							clusterID)

						status.Status = StatusDone
						status.Timestamp = getTimestamp()
						setStatus(clusterID, status, true)
					}
					break
				case StatusRunning:
					logger.GetInfo().Printf("monitor reported %s for cluster %s",
						status.Status, clusterID)
					break
				case StatusDone, StatusError:
					logger.GetInfo().Printf("monitor reported %s for cluster %s",
						status.Status, clusterID)
					if currentTime-status.Timestamp > doneReportTime {
						terminateCluster(client)
						status.Status = StatusTerminating
						status.Timestamp = getTimestamp()
						setStatus(clusterID, status, true)
					}
					break
				case StatusCanceled:
					logger.GetInfo().Printf("monitor reported %s for cluster %s",
						status.Status, clusterID)
					if currentTime-status.Timestamp > cancelTerminationDelay {
						terminateCluster(client)
						status.Status = StatusTerminating
						status.Timestamp = getTimestamp()
						setStatus(clusterID, status, true)
					}
					break
				case StatusTerminating:
					logger.GetInfo().Printf("monitor reported %s for cluster %s",
						status.Status, clusterID)
					if client.DestructionConfirmed() {
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
		releaseClusterLock(clusterID)
	}
}

func releaseLock() {
	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()

	redisClient.Del(monitorLock)
}

func acquireMonitorLock() bool {
	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()

	id, err := os.Hostname()
	if err != nil {
		logger.GetError().Println(err)
		return false
	}

	redisClient.SetNX(monitorLock, id, 15*time.Minute).Val()
	return id == redisClient.Get(monitorLock).Val()
}

func lockSuccess(clusterID string, prefix string, lockExpiration time.Duration) bool {
	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()
	key := clusterLockPreifx + clusterID
	value := prefix + strconv.FormatInt(time.Now().UnixNano(), 10)

	return redisClient.SetNX(key, value, lockExpiration*time.Second).Val()
}

func releaseClusterLock(clusterID string) {
	redisClient := datastore.GetRedisClient()
	defer redisClient.Close()
	key := clusterLockPreifx + clusterID
	redisClient.Del(key)
}

func acquireClusterLock(clusterID string, lockID string, lockExpiration time.Duration) error {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		if lockSuccess(clusterID, lockID, lockExpiration) {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("failed to acquire lock for cluster " +
		clusterID + ", lockID " + lockID)
}

func getTimestamp() int64 {
	return time.Now().Unix()
}
