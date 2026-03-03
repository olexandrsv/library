package log

import (
	"fmt"
	"library/errors"
	"net/http"
	"os"
)

var l *logger

type logger struct {
	logFile   LogFile
	traceFile LogFile
}

func Init() error {
	logFile, err := openFile("log.txt")
	if err != nil {
		return err
	}
	traceFile, err := openFile("trace.txt")
	if err != nil {
		return err
	}
	l = &logger{
		logFile:   newFile(logFile),
		traceFile: newFile(traceFile),
	}
	return nil
}

func openFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.NewFileErr(err, path, errors.NewInternalError())
	}
	return file, nil
}

func Info(msg string) {
	l.logFile.print(msg)
}

func Infof(pattern string, values ...any) {
	msg := fmt.Sprintf(pattern, values...)
	l.logFile.print(msg)
}

func Trace(msg string) {
	l.traceFile.print(msg)
}

func Error(msg string) {
	l.logFile.print(msg)
}

func Endpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, err := l.logFile.read()
	if err != nil {
		fmt.Printf("log.Endpoint l.read() error: %+v", err)
		return
	}

	if _, err := w.Write([]byte(body)); err != nil {
		fmt.Printf("log.Endpoint w.Write() error: %+v", err)
	}
}

func TraceEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, err := l.traceFile.read()
	if err != nil {
		fmt.Printf("log.TraceEndpoint l.traceFile.read() error: %+v", err)
		return
	}

	if _, err := w.Write([]byte(body)); err != nil {
		fmt.Printf("log.TraceEndpoint w.Write() error: %+v", err)
	}
}
