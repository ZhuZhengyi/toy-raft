// log.go

package raft

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Detail(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	SetLogLevel(logLevel string)
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelDetail
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

func parseLogLevel(logLevel string) LogLevel {
	switch strings.ToLower(logLevel) {
	case "detail":
		return LogLevelDetail
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	default:
		fmt.Fprintf(os.Stdout, "Parse log level err, unkown level(%v), set LogLevel with %v\n", logLevel, LogLevelInfo)
		return LogLevelInfo
	}
}

var (
	stdLogger = NewStdLogger(LogLevelInfo)
)

var (
	_ (Logger) = (*StdLogger)(nil)
)

type StdLogger struct {
	logLevel LogLevel
	info     *log.Logger
	debug    *log.Logger
	detail   *log.Logger
	warn     *log.Logger
	errorl   *log.Logger
	fatal    *log.Logger
}

func NewStdLogger(logLevel LogLevel) *StdLogger {
	return &StdLogger{
		logLevel: logLevel,
		detail:   log.New(os.Stdout, "DETAIL ", log.Ldate|log.Ltime|log.Lshortfile),
		debug:    log.New(os.Stdout, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile),
		info:     log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:     log.New(os.Stdout, "WARN ", log.Ldate|log.Ltime|log.Lshortfile),
		errorl:   log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
		fatal:    log.New(os.Stderr, "FATAL ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *StdLogger) SetLogLevel(logLevel string) {
	l.logLevel = parseLogLevel(logLevel)
}

func (l *StdLogger) Info(format string, v ...interface{}) {
	if l.logLevel > LogLevelInfo {
		return
	}
	l.info.Output(2, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Debug(format string, v ...interface{}) {
	if l.logLevel > LogLevelDebug {
		return
	}
	l.debug.Output(2, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Detail(format string, v ...interface{}) {
	if l.logLevel > LogLevelDetail {
		return
	}
	l.debug.Output(2, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Warn(format string, v ...interface{}) {
	if l.logLevel > LogLevelWarn {
		return
	}
	l.warn.Output(2, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Error(format string, v ...interface{}) {
	if l.logLevel > LogLevelError {
		return
	}
	l.errorl.Output(2, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Fatal(format string, v ...interface{}) {
	if l.logLevel > LogLevelFatal {
		return
	}
	l.fatal.Output(2, fmt.Sprintf(format, v...))

	panic(fmt.Sprintf(format, v...))
}
