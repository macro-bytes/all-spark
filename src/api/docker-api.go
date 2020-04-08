package api

import (
	"cloud"
	"errors"
	"io/ioutil"
	"logger"
	"monitor"
	"net/http"
	"util/serializer"
)

func validateDockerTemplate(template cloud.DockerEnvironment) error {
	if len(template.ClusterID) == 0 ||
		template.MemBytes < 10 ||
		template.NanoCpus < 10 ||
		template.WorkerNodes < 2 ||
		len(template.Image) == 0 {
		return errors.New("invalid template object")
	}

	return nil
}

func validateDockerFormBody(r *http.Request) (*cloud.DockerEnvironment, error) {
	err := validateRequest(r, "POST")
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	logger.GetInfo().Printf("Form body: %s", buffer)

	var template cloud.DockerEnvironment
	err = serializer.Deserialize(buffer, &template)
	if err != nil {
		return nil, err
	}

	err = validateDockerTemplate(template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func terminateDocker(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("http-request: /docker/terminate")
	terminate(w, r, cloud.Docker)
}

func createClusterDocker(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("http-request: /docker/create")
	client, err := validateDockerFormBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	serializedClient, err := serializer.Serialize(client)
	if err != nil {
		logger.GetError().Println(err)
	}

	err = monitor.RegisterCluster(client.ClusterID, cloud.Docker, serializedClient)
	if err != nil {
		logger.GetError().Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = client.CreateCluster()
	if err != nil {
		logger.GetError().Println(err.Error())
		monitor.DeregisterCluster(client.ClusterID)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully launched cluster"))
}

// InitDockerAPI - Initialize the Docker REST API
func InitDockerAPI() {
	http.HandleFunc("/docker/create", createClusterDocker)
	http.HandleFunc("/docker/terminate", terminateDocker)
}
