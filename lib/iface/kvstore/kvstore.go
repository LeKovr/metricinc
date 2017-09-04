/*
KV-store interface.

Интерфейс kvstore.Store используется для разделения логики хранения данных (например, boltdb) 
и кода, который это хранение использует (например, grpcapi).

*/
package kvstore

import (
	"lekovr/exam/counter/setup"
)

// Store holds used methods
type Store interface {
	// GetNumber returns number from store
	GetNumber() (*int64, error)

	// SetNumber stores number
	SetNumber(number *int64) error

	// GetSettings returns settings from store
	GetSettings() (*setup.Settings, error)

	// SetSettings stores settings
	SetSettings(sets *setup.Settings) error

	// Close used for closing store connection
	Close() error
}
