package utility

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Log manages log
type Log struct {
	queue  chan string
	finish chan bool
}

// NewLog creates new instance of Log
func NewLog(c *Config) (*Log, error) {
	path, err := c.Get("LogFileName")
	if err != nil {
		return nil, err
	}
	writer, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(writer)

	queue := make(chan string, 10000)
	finish := make(chan bool)

	go func() {
		for message := range queue {
			log.Println(message)
		}
		close(finish)
	}()

	return &Log{queue: queue, finish: finish}, nil
}

// Error stores ERROR level message into log
func (l *Log) Error(where string, message string, err error) {
	l.queue <- fmt.Sprintf("ERROR %s %s %s %v", time.Now().UTC().Format(time.RFC3339), where, message, err)
}

// Trace stores TRACE level message into log
func (l *Log) Trace(where string, message string) {
	l.withLevel("TRACE", where, message)
}

func (l *Log) withLevel(level string, where string, message string) {
	l.queue <- fmt.Sprintf("%s %s %s %s", level, time.Now().UTC().Format(time.RFC3339), where, message)
}

// Close properly ends Log instance lifecycle
func (l *Log) Close() {
	close(l.queue)
	<-l.finish
}
