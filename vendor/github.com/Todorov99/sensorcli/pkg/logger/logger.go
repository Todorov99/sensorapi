package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogrus(pkg string, writer ...io.Writer) *logrus.Logger {
	logFile, err := os.OpenFile(fmt.Sprintf("%s.log", pkg), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	writers := []io.Writer{
		logFile,
	}

	writers = append(writers, writer...)

	log := &logrus.Logger{
		Out:       io.MultiWriter(writers...),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{},
	}

	log.WithFields(logrus.Fields{
		"package": pkg,
	})

	return log
}
