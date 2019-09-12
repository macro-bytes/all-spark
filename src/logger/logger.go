package logger

import (
	"log"
	"os"
)

// GetInfo - logs informational messages to STDOUT
func GetInfo() *log.Logger {
	return log.New(os.Stdout, "INFO: ", log.Lshortfile)
}

// GetError - logs error message to STDERR
func GetError() *log.Logger {
	return log.New(os.Stderr, "ERROR: ", log.Lshortfile)
}

// GetFatal - logs fatal error message to STDERR and exits
func GetFatal() *log.Logger {
	return log.New(os.Stderr, "FATAL: ", log.Lshortfile)
}
