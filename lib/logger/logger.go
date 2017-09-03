// Package logger incapsulates logger logic
// This is a proxy for github.com/Sirupsen/logrus
// with some additions like
// functional options and logging to file
package logger

import (
	"github.com/Sirupsen/logrus"

	"lekovr/exam/lib/struct/logger"
)

// Log is a copy of logrus.Entry
type Log struct {
	*logrus.Entry
}

// Options is a program flags sample
type Config struct {
	LogLevel string `long:"log_level" description:"Log level [warn|info|debug]" default:"debug"`
}

// -----------------------------------------------------------------------------

// NewLogger creates a logger object
func NewLogger(cfg Config) (logger.Entry, error) {
	log := logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	log.Level = level

	l := logrus.NewEntry(log)
	return l, nil
}

// TODO: добавить хук с выводом строки в файле:
//pc, file, line, _ := runtime.Caller(2)
