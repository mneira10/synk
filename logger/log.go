package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// This is basically a workaround to have a global logging object in all
// packages. I don't like it but it works. Inspired by:
// https://stackoverflow.com/a/30261304/10296312

var logger = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: &logrus.TextFormatter{TimestampFormat: "15:04:05.00000", FullTimestamp: true},
	Level:     logrus.DebugLevel,
}

type Fields map[string]interface{}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return logger.WithFields(fields)
}

func SetLogLevel(level logrus.Level) {
	logger.SetLevel(level)
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
