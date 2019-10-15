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

func destroyCluster(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("destroyCluster")
	err := validateRequest(r, "POST")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	clusterID := r.PostFormValue("clusterID")
	if len(clusterID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("clusterID not specified"))
		return
	}

	clientBuffer, environment, err := monitor.GetClientData(clusterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to retrieve status for clusterID " + clusterID))
		return
	}

	client, err := cloud.Create(environment, clientBuffer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to establish allspark client with clusterID " + clusterID))
		return
	}

	status := monitor.GetLastKnownStatus(clusterID)
	if status == monitor.StatusPending {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("unable to destroy cluster " + clusterID +
			" when status is " + monitor.StatusPending))
		return
	}

	if status != monitor.StatusDone {
		err = client.DestroyCluster()
		if err != nil {
			logger.GetInfo().Printf(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("There was an error destroying clusterID " + clusterID))
			return
		}
	}

	monitor.DeregisterCluster(clusterID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully destroyed cluster"))
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

func healthCheck(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("healthCheck")
	err := validateRequest(r, "GET")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	return
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
	http.HandleFunc("/healthCheck", healthCheck)
	http.HandleFunc("/destroyCluster", destroyCluster)
	http.ListenAndServe(":32418", nil)
}
