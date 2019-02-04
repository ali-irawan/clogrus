package clogrus

import (
	"fmt"
	"github.com/nbs-go/clog"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// constants
const DefaultHostname = "unknown"
const SkipTrace = 1

// init get logger configuration from env and register ConsoleLogger to Logger
func init() {
	// Retrieve hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = DefaultHostname
	}
	// Retrieve logger level output
	var level int8
	levelStr, ok := os.LookupEnv(clog.EnvLogLevel)
	// If env not set, set to default level
	if !ok {
		level = clog.DefaultLevel
	} else {
		// Parse to int
		tmp, err := strconv.ParseInt(levelStr, 10, 8)
		// If level is unable to parse, set to default level
		if err != nil {
			level = clog.DefaultLevel
		} else {
			// Else, set level
			level = int8(tmp)
		}
	}
	// Convert level
	lv := getLevel(level)
	// Initiate writer
	logrus.SetLevel(lv)
	// Initiate logger instance
	log := ConsoleLogger{
		hostname: hostname,
		level:    level,
		writer:   logrus.StandardLogger(),
	}
	// Register logger
	clog.Register(&log)
}

func getLevel(level int8) logrus.Level {
	switch level {
	case clog.LevelError:
		return logrus.ErrorLevel
	case clog.LevelWarn:
		return logrus.WarnLevel
	case clog.LevelInfo:
		return logrus.InfoLevel
	case clog.LevelDebug:
		return logrus.DebugLevel
	default:
		panic("nbs-go/clogrus: unsupported logger level.")
	}
}

// ConsoleLogger is an implementation of Logger that prints output to console
type ConsoleLogger struct {
	hostname string
	level    int8
	writer   *logrus.Logger
}

func (l *ConsoleLogger) addFields() *logrus.Entry {
	return l.writer.
		WithField("hostname", l.hostname)
}

func (l *ConsoleLogger) Debug(msg string) {
	l.addFields().Debug(msg)
}

func (l *ConsoleLogger) Debugf(format string, args ...interface{}) {
	l.addFields().Debugf(format, args...)
}

func (l *ConsoleLogger) Info(msg string) {
	l.addFields().Info(msg)
}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	l.addFields().Infof(format, args...)
}

func (l *ConsoleLogger) Warn(msg string) {
	l.addFields().Warn(msg)
}

func (l *ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.addFields().Warnf(format, args...)
}

func (l *ConsoleLogger) Error(msg string, err error) {
	// Trace error
	file, line := clog.Trace(SkipTrace)
	l.addFields().
		WithField("trace", fmt.Sprintf("%s:%d", file, line)).
		WithField("error_msg", err).
		Error(msg)
}

func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	// Trace error
	file, line := clog.Trace(SkipTrace)
	l.addFields().
		WithField("trace", fmt.Sprintf("%s:%d", file, line)).
		Errorf(format, args...)
}
