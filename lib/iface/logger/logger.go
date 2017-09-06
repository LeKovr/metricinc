/*
Package logger defines log interface.

Интерфейс logger.Entry используется для разделения библиотеки журналирования (например, logrus)
и кода, который это журналирование использует (например, lib/grpcapi, lib/boltdb).

Интерфейс содержит сигнатуры стандартных методов журналирования и WithField, который у базовой библиотеки (logrus) возвращает внутренний тип, а не этот интерфейс. Поэтому для WithField нужна обертка (см lib/logger).
*/
package logger

// Entry is an interface which allows mocks
type Entry interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})

	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})

	WithField(key string, value interface{}) Entry
}
