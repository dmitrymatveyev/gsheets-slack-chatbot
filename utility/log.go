package utility

import (
	"log"
	"os"
	"time"
)

// Log manages log
type Log struct{}

// NewLog creates new instance of Log
func NewLog(c *Config) (*Log, error) {
	path, err := c.Get("logFileName")
	if err != nil {
		return nil, err
	}
	writer, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(writer)
	return &Log{}, nil
}

// Error stores ERROR level message into log
func (l *Log) Error(where string, message string, err error) {
	log.Printf("ERROR %s %s %s %v\n", time.Now().UTC().Format(time.RFC3339), where, message, err)
}

// Trace stores TRACE level message into log
func (l *Log) Trace(where string, message string) {
	withLevel("TRACE", where, message)
}

func withLevel(level string, where string, message string) {
	log.Printf("%s %s %s %s\n", level, time.Now().UTC().Format(time.RFC3339), where, message)
}
