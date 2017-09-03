package kvstore

import (
	"lekovr/exam/counter/setup"
)

type Store interface {
	// Works with KV store

	Close() error
	GetSettings() (*setup.Settings, error)
	SetSettings(sets *setup.Settings) error
	GetNumber() (*int64, error)
	SetNumber(number *int64) error
}
