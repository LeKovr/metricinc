/*
Package logger incapsulates logger logic.

This is a proxy for https://github.com/Sirupsen/logrus

TODO

Добавить хук с выводом строки в файле:
 pc, file, line, _ := runtime.Caller(2)

Варианты:

* https://github.com/sirupsen/logrus/issues/63

* https://github.com/prometheus/common/blob/master/log/log.go#L234

*/
package logger

import (
	"github.com/Sirupsen/logrus"
	"os"

	"github.com/LeKovr/metricinc/lib/iface/logger"
)

// Log is a copy of logrus.Entry
type Log struct {
	*logrus.Entry
}

// Config is a program flags group used in constructor
type Config struct {
	Level     string `long:"log_level" description:"Log level [warn|info|debug]" default:"debug"`
	UseStdOut bool   `long:"log_stdout" description:"Log to STDOUT without color and timestamps"`
}

// -----------------------------------------------------------------------------

// NewLogger creates a logger object
func NewLogger(cfg Config) (logger.Entry, error) {
	logDest := logrus.New()
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	logDest.Level = level
	if cfg.UseStdOut {
		logDest.Out = os.Stdout
		logDest.Formatter = &logrus.TextFormatter{DisableColors: true, DisableTimestamp: true}
	}
	logEntry := logrus.NewEntry(logDest)

	log := Log{logEntry}
	log.WithField("config", cfg).Debug("Create logger")

	return &log, nil

}

// WithField call same name parent method and returns result as logger.Entry
func (entry Log) WithField(key string, value interface{}) logger.Entry {
	l := entry.Entry.WithField(key, value)
	return &Log{l}
}
