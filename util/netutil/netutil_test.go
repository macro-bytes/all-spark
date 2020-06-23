package netutil

import (
	"strconv"
	"testing"
	"time"
)

func TestIsListeningOnPort(t *testing.T) {
	closedPort := 60713

	/*
		if !IsListeningOnPort("localhost", openPort, 10*time.Second, 1) {
			t.Error("Expected host to be listening on port " +
				strconv.Itoa(openPort))
		}
	*/

	if IsListeningOnPort("localhost", closedPort, 10*time.Second, 1) {
		t.Error("Expected host not to be listening on port " +
			strconv.Itoa(closedPort))
	}
}
