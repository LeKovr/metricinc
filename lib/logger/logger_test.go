package logger_test

import (
	"lekovr/exam/lib/logger"
)

// Example_usage is a usage example
func Example_usage() {

	// Prepare config
	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "warn"
	cfg.Logger.UseStdOut = true

	// Create object
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}
	// Do logging
	log.Warn("Start server")
	// Output: level=warning msg="Start server"
}
