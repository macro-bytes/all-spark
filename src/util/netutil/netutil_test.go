package netutil

import (
	"testing"
	"time"
)

func TestIsListeningOnPort(t *testing.T) {
	openPort := "22"
	closedPort := "60713"

	if !IsListeningOnPort("localhost", openPort, 10*time.Second) {
		t.Error("Expected host to be listening on port " +
			openPort)
	}

	if IsListeningOnPort("localhost", closedPort, 10*time.Second) {
		t.Error("Expected host not to be listening on port " +
			closedPort)
	}
}
