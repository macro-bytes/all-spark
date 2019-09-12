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

func validateAwsTemplate(template cloud.AwsEnvironment) error {
	if len(template.ClusterID) == 0 ||
		template.EBSVolumeSize < 10 ||
		len(template.IAMRole) == 0 ||
		len(template.ImageID) == 0 ||
		len(template.InstanceType) == 0 ||
		len(template.Region) == 0 ||
		len(template.SecurityGroupIds) == 0 ||
		len(template.SubnetID) == 0 ||
		template.WorkerNodes < 2 {
		return errors.New("invalid template object")
	}

	return nil
}

func validateAwsFormBody(r *http.Request) (*cloud.AwsEnvironment, error) {
	err := validateRequest(r, "POST")
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var template cloud.AwsEnvironment

	err = decoder.Decode(&template)
	if err != nil {
		return nil, err
	}

	err = validateAwsTemplate(template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func createClusterAws(w http.ResponseWriter, r *http.Request) {
	client, err := validateAwsFormBody(r)
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

	monitor.RegisterCluster(client.ClusterID, cloud.Aws, serializedClient)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully launched cluster"))
}

func destroyClusterAws(w http.ResponseWriter, r *http.Request) {
	client, err := validateAwsFormBody(r)
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

// InitAwsAPI - Initialize the AWS REST API
func InitAwsAPI() {
	http.HandleFunc("/aws/createCluster", createClusterAws)
	http.HandleFunc("/aws/destroyCluster", destroyClusterAws)
}
