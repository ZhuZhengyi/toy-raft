// log.go

package raft

import (
	"log"
	"os"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

var (
	logger = NewStdLogger()
)

type StdLogger struct {
	info   *log.Logger
	warn   *log.Logger
	errorl *log.Logger
	debug  *log.Logger
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		info:   log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:   log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorl: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debug:  log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *StdLogger) Info(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func (l *StdLogger) Warn(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}
func (l *StdLogger) Error(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func (l *StdLogger) Debug(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}
