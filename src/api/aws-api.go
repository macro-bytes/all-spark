package api

import (
	"cloud"
	"encoding/json"
	"errors"
	"net/http"
	"template"
)

func validateAwsTemplate(template template.AwsTemplate) error {
	if len(template.ClusterID) == 0 ||
		template.EBSVolumeSize < 10 ||
		len(template.IAMRole) == 0 ||
		len(template.ImageId) == 0 ||
		len(template.InstanceType) == 0 ||
		len(template.Region) == 0 ||
		len(template.SecurityGroupIds) == 0 ||
		len(template.SubnetId) == 0 ||
		template.WorkerNodes < 2 {
		return errors.New("invalid template object")
	}

	return nil
}

func createClusterAws(w http.ResponseWriter, r *http.Request) {
	err := validatePostRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var template template.AwsTemplate

	err = decoder.Decode(&template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = validateAwsTemplate(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	client := &cloud.AwsEnvironment{}
	_, err = client.CreateClusterHelper(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully launched cluster"))
}

func destroyClusterAws(w http.ResponseWriter, r *http.Request) {
	err := validatePostRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var template template.AwsTemplate

	err = decoder.Decode(&template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = validateAwsTemplate(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	client := &cloud.AwsEnvironment{}
	err = client.DestroyClusterHelper(template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully destroyed cluster"))
}

// InitAwsAPI - Initialize the AWS REST API
func InitAwsAPI() {
	http.HandleFunc("/aws/createCluster", createClusterAws)
	http.HandleFunc("/aws/destroyCluster", destroyClusterAws)
	http.HandleFunc("/aws/checkIn", checkIn)
	http.ListenAndServe(":32418", nil)
}
