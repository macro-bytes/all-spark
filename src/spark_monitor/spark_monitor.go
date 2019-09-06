package spark_monitor

import (
	"time"
)

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

func HandleCheckIn(clusterID, clusterStatus []byte) {

}

func SetIdle(clusterID string, clusterStatus []byte) {

}

func SetPending(clusterID string, clusterStatus []byte) {

}

func SetRunning(clusterID string, clusterStatus []byte) {

}

func RunClusterMonitor(iterations int) {
	if iterations <= 0 {
		for {
			runClusterMonitorHelper()
			time.Sleep(10 * time.Second)
		}
	}
	for i := 0; i < iterations; i++ {
		runClusterMonitorHelper()
		time.Sleep(10 * time.Second)
	}
}

func runClusterMonitorHelper() {

}
