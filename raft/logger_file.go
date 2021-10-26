// log.go

package raft

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	_ (Logger) = (*FileLogger)(nil)
)

//FileLogger
type FileLogger struct {
	file     *os.File
	logLevel LogLevel
	info     *log.Logger
	debug    *log.Logger
	warn     *log.Logger
	errorl   *log.Logger
	fatal    *log.Logger
}

func NewFileLogger(filePath string) *FileLogger {
	// open file
	dirPath := filepath.Dir(filePath)
	if dir, err := os.Stat(dirPath); err != nil {
		os.MkdirAll(dirPath, 0644)
	} else if !dir.IsDir() {
		fmt.Printf("dirPath: %v is not directory\n", dirPath)
		return nil
	}

	//os.OpenFile(filePath, "", 0644)

	return &FileLogger{
		debug:  log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		info:   log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:   log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorl: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatal:  log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *FileLogger) SetLogLevel(logLevel string) {
	l.logLevel = parseLogLevel(logLevel)
}

func (l *FileLogger) Info(format string, v ...interface{}) {
	if l.logLevel > LogLevelInfo {
		return
	}
	l.info.Printf(format, v...)
}

func (l *FileLogger) Debug(format string, v ...interface{}) {
	if l.logLevel > LogLevelDebug {
		return
	}
	l.debug.Printf(format, v...)
}

func (l *FileLogger) Warn(format string, v ...interface{}) {
	if l.logLevel > LogLevelWarn {
		return
	}
	l.warn.Printf(format, v...)
}

func (l *FileLogger) Error(format string, v ...interface{}) {
	if l.logLevel > LogLevelError {
		return
	}
	l.errorl.Printf(format, v...)
}

func (l *FileLogger) Fatal(format string, v ...interface{}) {
	if l.logLevel > LogLevelFatal {
		return
	}
	l.fatal.Printf(format, v...)

	panic(fmt.Sprintf(format, v...))
}
