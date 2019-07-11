package netutil

import (
	"net"
	"time"
)

func IsListeningOnPort(host string, port string, timeout time.Duration) bool {
	retries := 3

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
		time.Sleep(5 * time.Second)
	}

	return false
}
