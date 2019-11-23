package logger

import (
	"log"
	"os"
	"time"
)

// GetInfo - logs informational messages to STDOUT
func GetInfo() *log.Logger {
	return log.New(os.Stdout,
		time.Now().Format("[2006-01-02 15:04:05]")+" INFO: ",
		log.Lshortfile)
}

// GetDebug - logs informational messages to STDOUT
func GetDebug() *log.Logger {
	return log.New(os.Stdout,
		time.Now().Format("[2006-01-02 15:04:05]")+" DEBUG: ",
		log.Lshortfile)
}

// GetError - logs error message to STDERR
func GetError() *log.Logger {
	return log.New(os.Stderr,
		time.Now().Format("[2006-01-02 15:04:05]")+" ERROR: ",
		log.Lshortfile)
}

// GetFatal - logs fatal error message to STDERR and exits
func GetFatal() *log.Logger {
	return log.New(os.Stderr,
		time.Now().Format("[2006-01-02 15:04:05]")+" FATAL: ",
		log.Lshortfile)
}
