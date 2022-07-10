package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00000", FullTimestamp: true},
	Level:     logrus.DebugLevel,
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Panic(args ...interface{}) {
	logger.Panic(args...)
}
