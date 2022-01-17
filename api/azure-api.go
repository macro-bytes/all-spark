package api

import (
	"allspark/cloud"
	"allspark/logger"
	"allspark/monitor"
	"allspark/util/serializer"
	"errors"
	"io/ioutil"
	"net/http"
)

func validateAzureTemplate(template cloud.AzureEnvironment) error {
	if len(template.ClusterID) == 0 ||
		len(template.SubscriptionID) == 0 ||
		len(template.Region) == 0 ||
		len(template.ClientID) == 0 ||
		len(template.ClientSecret) == 0 ||
		len(template.Tenant) == 0 ||
		len(template.ResourceGroup) == 0 ||
		len(template.VMNet) == 0 ||
		len(template.VMSubnet) == 0 ||
		len(template.VMSize) == 0 ||
		len(template.ImageStorageAccount) == 0 ||
		len(template.ImageContainer) == 0 ||
		len(template.ImageBlob) == 0 ||
		template.WorkerNodes < 0 {
		return errors.New("invalid template object")
	}

	return nil
}

func validateAzureFormBody(r *http.Request) (*cloud.AzureEnvironment, error) {
	err := validateRequest(r, "POST")
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var template cloud.AzureEnvironment
	err = serializer.Deserialize(buffer, &template)
	if err != nil {
		return nil, err
	}

	err = validateAzureTemplate(template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func terminateAzure(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("http-request: /azure/terminate")
	terminate(w, r, cloud.Azure)
}

func createClusterAzure(w http.ResponseWriter, r *http.Request) {
	logger.GetInfo().Println("http-request: /azure/create")
	client, err := validateAzureFormBody(r)
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

	logger.GetInfo().Println("http-request: /azure/create, clusterID: " + client.ClusterID)

	err = monitor.RegisterCluster(client.ClusterID, cloud.Azure, serializedClient)
	if err != nil {
		logger.GetError().Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = client.CreateCluster()
	if err != nil {
		monitor.SetCanceled(client.ClusterID)
		logger.GetError().Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully launched cluster"))
}

// InitAzureAPI - Initialize the Azure API
func InitAzureAPI() {
	http.HandleFunc("/azure/create", createClusterAzure)
	http.HandleFunc("/azure/terminate", terminateAzure)
}
