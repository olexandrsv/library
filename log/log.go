package log

import (
	"fmt"
	"net/http"
)

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Trace(string)
	Error(error)
	ReadLogs() (string, error)
	ReadTraces() (string, error)
}

var l Logger

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

func InitMock(mock MockLog) {
	l = &mock
}

func (l *logger) Info(msg string) {
	l.logFile.print(msg)
	fmt.Println(msg)
}

func (l *logger) Infof(pattern string, values ...any) {
	msg := fmt.Sprintf(pattern, values...)
	l.logFile.print(msg)
	fmt.Println(msg)
}

func (l *logger) Trace(msg string) {
	l.traceFile.print(msg)
	fmt.Println(msg)
}

func (l *logger) Error(err error) {
	if err != nil {
		l.logFile.print(err.Error())
		fmt.Println(err.Error())
	}
}

func (l *logger) ReadLogs() (string, error) {
	return l.logFile.read()
}

func (l *logger) ReadTraces() (string, error) {
	return l.traceFile.read()
}

func Info(msg string) {
	l.Info(msg)
}

func Infof(pattern string, values ...any) {
	l.Infof(pattern, values...)
}

func Trace(msg string) {
	l.Trace(msg)
}

func Error(err error) {
	l.Error(err)
}

func Endpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, err := l.ReadLogs()
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

	body, err := l.ReadTraces()
	if err != nil {
		fmt.Printf("log.TraceEndpoint l.traceFile.read() error: %+v", err)
		return
	}

	if _, err := w.Write([]byte(body)); err != nil {
		fmt.Printf("log.TraceEndpoint w.Write() error: %+v", err)
	}
}
