package api

import (
	"net/http"
)

func createClusterDocker(w http.ResponseWriter, r *http.Request) {

}

func destroyClusterDocker(w http.ResponseWriter, r *http.Request) {

}

// InitDockerAPI - Initialize the Docker REST API
func InitDockerAPI() {
	http.HandleFunc("/aws/createCluster", createClusterAws)
	http.HandleFunc("/aws/destroyCluster", destroyClusterAws)
	http.HandleFunc("/aws/checkIn", checkIn)
	http.ListenAndServe(":32418", nil)
}
