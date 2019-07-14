package netutil

import (
	"net"
	"time"
)

func IsListeningOnPort(host string, port string,
	timeout time.Duration, retries int) bool {

	for retries > 0 {
		conn, err := net.DialTimeout("tcp",
			net.JoinHostPort(host, port), timeout)
		if conn != nil {
			conn.Close()
		}

		if err == nil {
			return true
		}

		retries--
		time.Sleep(1 * time.Second)
	}

	return false
}
