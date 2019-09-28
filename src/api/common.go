package api

import (
	"cloud"
	"daemon"
	"errors"
	"io/ioutil"
	"logger"
	"monitor"
	"net/http"
	"util/serializer"
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
	logger.GetInfo().Println("getStatus")
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
	logger.GetInfo().Printf("checking status on clusterID %v", clusterID)

	status := monitor.GetLastKnownStatus(clusterID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}

func checkIn(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("checkin")
	err := validateRequest(r, "POST")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var body cloud.SparkStatusCheckIn
	buffer, err := ioutil.ReadAll(r.Body)
	serializer.Deserialize(buffer, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	logger.GetInfo().Printf("Form body: %s", buffer)

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
		logger.GetFatal().Fatalln("invalid cloud environment specified")
	}

	http.HandleFunc("/checkIn", checkIn)
	http.HandleFunc("/getStatus", getStatus)
	http.ListenAndServe(":32418", nil)
}
