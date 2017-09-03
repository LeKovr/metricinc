package logger

// Logger is an interface which allows mocks
type Entry interface {
	Fatalf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})

	Info(args ...interface{})
	//	WithField(key string, value interface{}) *Entry
}
