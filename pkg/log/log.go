// log.go

package log

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(logPath string, logLevel string) {
	logDir := path.Dir(logPath)
	if err := os.MkdirAll(logDir, 0666); err != nil {
		log.Fatalf("Mkdir log file %v err: %v", logDir, err)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Open log file %v err: %v", logPath, err)
	}

	if strings.ToLower(logLevel) == "debug" {
		Debug = log.New(io.MultiWriter(logFile, os.Stdout), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Debug = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	Info = log.New(io.MultiWriter(logFile, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(logFile, os.Stdout), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(logFile, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
