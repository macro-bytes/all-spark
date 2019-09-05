package api

import (
	"cloud"
	"encoding/json"
	"errors"
	"net/http"
	"template"
)

func validateDockerTemplate(template template.DockerTemplate) error {
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

func createClusterDocker(w http.ResponseWriter, r *http.Request) {
	err := validatePostRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var template template.DockerTemplate

	err = decoder.Decode(&template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = validateDockerTemplate(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	client := &cloud.DockerEnvironment{}
	_, err = client.CreateClusterHelper(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully launched cluster"))
}

func destroyClusterDocker(w http.ResponseWriter, r *http.Request) {
	err := validatePostRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var template template.DockerTemplate

	err = decoder.Decode(&template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = validateDockerTemplate(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	client := &cloud.DockerEnvironment{}
	err = client.DestroyClusterHelper(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully destroyed cluster"))
}

// InitDockerAPI - Initialize the Docker REST API
func InitDockerAPI() {
	http.HandleFunc("/docker/createCluster", createClusterAws)
	http.HandleFunc("/docker/destroyCluster", destroyClusterAws)
	http.HandleFunc("/docker/checkIn", checkIn)
	http.ListenAndServe(":32418", nil)
}
