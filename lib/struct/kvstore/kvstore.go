package kvstore

import (
	"lekovr/exam/lib/struct/server"
)

type Store interface {
	// Works with KV store

	//	Open(cfg interface{}, file string) (Store, error)
	Close() error
	GetSettings() (*server.Settings, error)
	SetSettings(sets *server.Settings) error
	GetNumber() (*int64, error)
	SetNumber(number *int64) error
}

var (
// ErrDatabaseNotOpen is returned when a DB instance is accessed before it
// is opened or after it is closed.
//ErrSettingsEmpty = errors.New("settingsdatabase not open")
)
