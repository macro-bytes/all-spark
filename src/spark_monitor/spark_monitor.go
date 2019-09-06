package spark_monitor

import (
	"datastore"
	"strconv"
	"time"
)

const (
	statusPending = "STATUS_PENDING"
	statusIdle    = "STATUS_IDLE"
	statusRunning = "STATUS_RUNNING"
)

// SparkWorker describes the spark worker node state
type SparkWorker struct {
	ID            string `json:"id"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	WebUIAddress  string `json:"webuiaddress"`
	Cores         int    `json:"cores"`
	CoresUsed     int    `json:"coresused"`
	CoresFree     int    `json:"coresfree"`
	Memory        uint64 `json:"memory"`
	MemoryUsed    uint64 `json:"memoryused"`
	MemoryFree    uint64 `json:"memoryfree"`
	State         string `json:"state"`
	LastHeartBeat uint64 `json:"lastheartbeat"`
}

// SparkApp describes the spark application state
type SparkApp struct {
	ID             string `json:"id"`
	StartTime      uint64 `json:"starttime"`
	Name           string `json:"name"`
	Cores          int    `json:"cores"`
	User           string `json:"user"`
	MemoryPerSlave int    `json:"memoryperslave"`
	SubmitDate     string `json:"submitdate"`
	State          string `json:"state"`
	Duration       uint64 `json:"duration"`
}

// SparkClusterStatus describes the entire spark cluster state
type SparkClusterStatus struct {
	URL           string        `json:"url"`
	Workers       []SparkWorker `json:"workers"`
	AliveWorkers  int           `json:"aliveworkers"`
	Cores         int           `json:"cores"`
	CoresUsed     int           `json:"coresused"`
	Memory        uint64        `json:"memory"`
	MemoryUsed    uint64        `json:"memoryused"`
	ActiveApps    []SparkApp    `json:"activeapps"`
	CompletedApps []SparkApp    `json:"completedapps"`
	Status        string        `json:"status"`
}

// SparkClusterStatusAtEpoch describes the state of a cluster
// at a given timestamp
type SparkClusterStatusAtEpoch struct {
	ClusterStatus SparkClusterStatus
	Timestamp     int64
}

// HandleCheckIn - handles spark monitor check-in http requests
func HandleCheckIn(clusterID, clusterStatus []byte) {

}

// SetPending - sets spark cluster status to pending
func SetPending(clusterID string, clusterStatus []byte) {
	register(statusPending, clusterID, clusterStatus)
}

// SetIdle - sets spark cluster status to idle
func SetIdle(clusterID string, clusterStatus []byte) {
	register(statusIdle, clusterID, clusterStatus)
}

// SetRunning - sets spark cluster status to running (i.e. job is running)
func SetRunning(clusterID string, clusterStatus []byte) {
	register(statusRunning, clusterID, clusterStatus)
}

func register(hmap string, clusterID string, status []byte) {
	client := datastore.GetRedisClient()
	defer client.Close()

	client.HSet(hmap, clusterID, string(status))
}

// MonitorSparkClusters - daemon used for monitoring all spark clusters;
// monitor will run for the specified number of iterations, or indefinitely
// if iterations <= 0.
func MonitorSparkClusters(iterations int) {
	if iterations <= 0 {
		for {
			monitorClusterHelper()
			time.Sleep(10 * time.Second)
		}
	}
	for i := 0; i < iterations; i++ {
		monitorClusterHelper()
		time.Sleep(10 * time.Second)
	}
}

func monitorClusterHelper() {

}

func getTimestamp() int64 {
	return time.Now().Unix()
}

func getTimestampAsString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
