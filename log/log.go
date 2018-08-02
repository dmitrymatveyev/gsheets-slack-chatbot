package log

import (
	"log"
	"time"
)

func Error(where string, message string, err error) {
	log.Printf("ERROR %s %s %s %v\n", time.Now().UTC().Format(time.RFC3339), where, message, err)
}

func Trace(where string, message string) {
	withLevel("TRACE", where, message)
}

func withLevel(level string, where string, message string) {
	log.Printf("%s %s %s %s\n", level, time.Now().UTC().Format(time.RFC3339), where, message)
}
