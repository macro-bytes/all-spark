package api

import (
	"errors"
	"net/http"
)

type sparkWorker struct {
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

type sparkApp struct {
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

func validatePostRequest(r *http.Request) error {
	if r.Method != "POST" {
		return errors.New("invalid request method: " + r.Method)
	}

	if r.Body == nil {
		return errors.New("form body is null")
	}

	return nil
}

func checkIn(w http.ResponseWriter, r *http.Request) {

}
