package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Writer() io.Writer
}

type logger struct {
	*logrus.Logger
}

func New() Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	l.SetLevel(logrus.DebugLevel)

	return &logger{l}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *logger) Writer() io.Writer {
	return l.Logger.Out
}
