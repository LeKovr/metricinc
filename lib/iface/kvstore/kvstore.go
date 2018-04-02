/*
Package kvstore defines Store interface.

Интерфейс kvstore.Store используется для разделения логики хранения данных (например, boltdb)
и кода, который это хранение использует (например, grpcapi).

*/
package kvstore

//go:generate mockgen -destination=../../mocks/kvstore.go -package=mocks github.com/LeKovr/metricinc/lib/iface/kvstore Store

import (
	"github.com/LeKovr/metricinc/counter/setup"
)

// Store holds used methods
type Store interface {

	// GetSettings returns settings from store
	GetSettings() (*setup.Settings, error)

	// SetSettings stores settings
	SetSettings(sets *setup.Settings) error

	// GetNumber returns number from store
	GetNumber() (*int64, error)

	// SetNumber stores number
	SetNumber(number *int64) error

	// Close used for closing store connection
	Close() error
}
