package logger_test

import (
	"lekovr/exam/lib/logger"
)

// This is a usage example:
func Example_usage() {

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "warn"
	cfg.Logger.UseStdOut = true
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}
	log.Warn("Start server")
	// Output: level=warning msg="Start server"
}
