package api

import (
	"cloud"
	"encoding/json"
	"errors"
	"monitor"
	"net/http"
)

func validatePostRequest(r *http.Request) error {
	if r.Method != "POST" {
		return errors.New("invalid request method: " + r.Method)
	}

	if r.Body == nil {
		return errors.New("form body is null")
	}

	return nil
}

func checkIn(w http.ResponseWriter, r *http.Request) {
	err := validatePostRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var body cloud.SparkStatusCheckIn
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	monitor.HandleCheckIn(body.ClusterID, body.Status)
}
