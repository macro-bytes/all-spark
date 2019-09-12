package api

import (
	"cloud"
	"encoding/json"
	"errors"
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
		len(template.Image) == 0 ||
		len(template.Network) == 0 {
		return errors.New("invalid template object")
	}

	return nil
}

func validateDockerFormBody(r *http.Request) (*cloud.DockerEnvironment, error) {
	err := validateRequest(r, "POST")
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var template cloud.DockerEnvironment

	err = decoder.Decode(&template)
	if err != nil {
		return nil, err
	}

	err = validateDockerTemplate(template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func createClusterDocker(w http.ResponseWriter, r *http.Request) {
	client, err := validateDockerFormBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = client.CreateCluster()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	serializedClient, err := serializer.Serialize(client)
	if err != nil {
		logger.Error(err.Error())
	}

	monitor.RegisterCluster(client.ClusterID, cloud.Docker, serializedClient)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully launched cluster"))
}

func destroyClusterDocker(w http.ResponseWriter, r *http.Request) {
	client, err := validateDockerFormBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = client.DestroyCluster()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	monitor.DeregisterCluster(client.ClusterID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully destroyed cluster"))
}

// InitDockerAPI - Initialize the Docker REST API
func InitDockerAPI() {
	http.HandleFunc("/docker/createCluster", createClusterDocker)
	http.HandleFunc("/docker/destroyCluster", destroyClusterDocker)
}
