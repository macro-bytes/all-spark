package logger

import (
	"log"
	"os"
)

// Info - logs informational messages to STDOUT
func Info(format string, v ...interface{}) {
	l := log.New(os.Stdout, "INFO: ", log.Llongfile)
	l.Printf(format + "\n")
}

// Error - logs error message to STDERR
func Error(format string, v ...interface{}) {
	l := log.New(os.Stderr, "ERROR: ", log.Llongfile)
	l.Printf(format + "\n")
}

// Fatal - logs fatal error message to STDERR and exits
func Fatal(format string, v ...interface{}) {
	l := log.New(os.Stderr, "FATAL: ", log.Llongfile)
	l.Fatalf(format + "\n")
}
