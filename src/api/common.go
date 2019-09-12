package api

import (
	"cloud"
	"daemon"
	"encoding/json"
	"errors"
	"logger"
	"monitor"
	"net/http"
)

func validateRequest(r *http.Request, method string) error {
	if r.Method != method {
		return errors.New("invalid request method: " + r.Method)
	}

	if method == "POST" && r.Body == nil {
		return errors.New("form body is null")
	}

	return nil
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	err := validateRequest(r, "GET")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	clusterID := r.FormValue("clusterID")
	if len(clusterID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("clusterID not specified"))
		return
	}

	status := monitor.GetLastKnownStatus(clusterID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}

func checkIn(w http.ResponseWriter, r *http.Request) {
	err := validateRequest(r, "POST")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var body cloud.SparkStatusCheckIn
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	monitor.HandleCheckIn(body.ClusterID, body.Status)
}

// Init - initializes the allspark-orchestrator web api
func Init() {
	switch daemon.GetAllSparkConfig().CloudEnvironment {
	case cloud.Aws:
		InitAwsAPI()
	case cloud.Docker:
		InitDockerAPI()
	default:
		logger.Fatal("invalid cloud environment specified")
	}

	http.HandleFunc("/checkIn", checkIn)
	http.HandleFunc("/getStatus", getStatus)
	http.ListenAndServe(":32418", nil)
}
