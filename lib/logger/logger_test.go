package logger_test

import (
	"lekovr/exam/lib/logger"
	"testing"
)

// Example_usage is an usage example
func Example_usage() {

	// Prepare config
	cfg := struct {
		Logger logger.Config
	}{
		logger.Config{Level: "warn", UseStdOut: true},
	}

	// Create object
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}
	// Do logging
	log.Warn("Start server")
	// Output: level=warning msg="Start server"
}

func TestLogger_NewLogger(t *testing.T) {

	// Prepare config
	cfg := struct {
		Logger logger.Config
	}{
		logger.Config{Level: "warn", UseStdOut: true},
	}

	// Create object
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		t.Errorf("%q. BadLogLevel() error: %+v", "good config", err)
	}
	log.Debug("skipped line")

}

func TestLogger_BadLogLevel(t *testing.T) {

	// Prepare config
	cfg := struct {
		Logger logger.Config
	}{
		logger.Config{Level: "unknown", UseStdOut: true},
	}

	// Create object
	log, err := logger.NewLogger(cfg.Logger)
	if err == nil {
		t.Errorf("%q. BadLogLevel() = %v, want error", "Bad log level", cfg)
		log.Debug("unreachable point")
	}

}
