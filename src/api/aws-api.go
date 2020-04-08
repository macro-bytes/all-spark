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

func validateAwsTemplate(template cloud.AwsEnvironment) error {
	if len(template.ClusterID) == 0 ||
		template.EBSVolumeSize < 10 ||
		len(template.IAMRole) == 0 ||
		len(template.Image) == 0 ||
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

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	logger.GetInfo().Printf("Form body: %s", buffer)

	var template cloud.AwsEnvironment
	err = serializer.Deserialize(buffer, &template)
	if err != nil {
		return nil, err
	}

	err = validateAwsTemplate(template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func terminateAws(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("http-request: /aws/terminate")
	terminate(w, r, cloud.Aws)
}

func createClusterAws(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("http-request: /aws/create")
	client, err := validateAwsFormBody(r)
	if err != nil {
		logger.GetError().Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	serializedClient, err := serializer.Serialize(client)
	if err != nil {
		logger.GetError().Println(err)
	}

	err = monitor.RegisterCluster(client.ClusterID, cloud.Aws, serializedClient)
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

// InitAwsAPI - Initialize the AWS REST API
func InitAwsAPI() {
	http.HandleFunc("/aws/create", createClusterAws)
	http.HandleFunc("/aws/terminate", terminateAws)
}
