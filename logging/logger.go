package logging

import (
	"log"
)

// Logger is a minimal logging interface used by the platform.
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// StdLogger implements Logger using the standard library.
type StdLogger struct{}

func (l StdLogger) Info(msg string)  { log.Println("INFO:", msg) }
func (l StdLogger) Error(msg string) { log.Println("ERROR:", msg) }
