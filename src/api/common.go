package api

import (
	"errors"
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

}
