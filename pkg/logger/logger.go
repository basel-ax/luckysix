// Package logger -.
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Interface -.
type Interface interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// Logger -.
type Logger struct {
	*logrus.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	l := logrus.New()

	l.SetReportCaller(true)
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	err := l.Level.UnmarshalText([]byte(level))
	if err != nil {
		l.Fatal(err)
	}

	return &Logger{l}
}
