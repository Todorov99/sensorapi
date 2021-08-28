package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	errorLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
}

func NewLogger(pkg string) Logger {
	logFile, err := os.OpenFile(fmt.Sprintf("%s.log", pkg), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		panic(err)
	}

	return Logger{
		errorLogger: log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(logFile, "INFO: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(logFile, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l Logger) Infof(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

func (l Logger) Warn(v ...interface{}) {
	l.warnLogger.Println(v...)
}

func (l Logger) Warnf(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}

func (l Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

func (l Logger) Errorf(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

func (l Logger) Panic(v interface{}) {
	l.errorLogger.Panic(v)
}
