package spark_monitor

import (
	"datastore"
	"log"
	"strconv"
	"time"
	"util/serializer"
)

const (
	statusPending = "PENDING"
	statusIdle    = "IDLE"
	statusRunning = "RUNNING"
	statusMap     = "STATUS_MAP"
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
	Timestamp int64
	Status    string
}

// HandleCheckIn - handles spark monitor check-in http requests
func HandleCheckIn(clusterID string, clusterStatus []byte) {
	log.Printf("cluster: %s, status: %s", clusterID, string(clusterStatus))
	var reportedStatus SparkClusterStatus
	serializer.Deserialize(clusterStatus, &reportedStatus)

	epochStatus := SparkClusterStatusAtEpoch{
		Timestamp: getTimestamp(),
		Status:    getReportedStatus(reportedStatus),
	}

	setStatus(clusterID, epochStatus)
}

// RegisterCluster - registers newly created spark
// cluster with a pending status
func RegisterCluster(clusterID string) {
	setStatus(clusterID, SparkClusterStatusAtEpoch{
		Status: statusPending,
	})
}

func getReportedStatus(status SparkClusterStatus) string {
	if len(status.ActiveApps) > 0 {
		return statusRunning
	}

	return statusIdle
}

func getPriorStatus(clusterID string) string {
	client := datastore.GetRedisClient()
	defer client.Close()

	var clusterState SparkClusterStatusAtEpoch
	err := serializer.Deserialize([]byte(client.HGet(statusMap, clusterID).Val()),
		&clusterState)
	if err != nil {
		log.Println(err)
	}

	return clusterState.Status
}

func setStatus(clusterID string, status SparkClusterStatusAtEpoch) {
	client := datastore.GetRedisClient()
	defer client.Close()

	result, err := serializer.Serialize(status)
	if err != nil {
		log.Println(err)
	}

	client.HSet(statusMap, clusterID, string(result))
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
