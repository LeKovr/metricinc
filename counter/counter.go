/*
Package counter contains counter logic.

Это основная математика проекта, поэтому размещена в корне, а не в каталоге внутренних библиотек lib/.
*/
package counter

// TODO: go:generate mockgen -destination=../lib/mocks/counter.go -package=mocks github.com/LeKovr/metricinc/counter Counter

import (
	"errors"
	"sync"

	"github.com/LeKovr/metricinc/counter/setup"
)

// Counter holds object internals
type Counter struct {
	number   int64
	settings *setup.Settings // nil means NewCounter was not called
	mutex    sync.RWMutex    // not referense, so defined always
}

// NewCounter creates a counter object
func NewCounter(s *setup.Settings, number int64) (*Counter, error) {

	if err := checkSettings(s); err != nil {
		return nil, err
	}
	cnt := Counter{number: number, settings: s}
	return &cnt, nil
}

// GetSettings returns current settings
func (c *Counter) GetSettings() (*setup.Settings, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.settings == nil {
		return nil, errors.New("NewCounter must be called before GetSettings")
	}
	return c.settings, nil
}

// SetSettings stores new settings
func (c *Counter) SetSettings(s *setup.Settings) error {
	if err := checkSettings(s); err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.settings = s
	return nil
}

// GetNumber returns current counter number
func (c *Counter) GetNumber() (*int64, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.settings == nil {
		return nil, errors.New("NewCounter must be called before GetNumber")
	}
	return &c.number, nil
}

// IncrementNumber adds c.settings.Step to current counter number.
// If new number >= c.settings.Limit, this limit deducted from number
func (c *Counter) IncrementNumber() (*int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.settings == nil {
		return nil, errors.New("NewCounter must be called before IncrementNumber")
	}

	// number + step == limit ? number = 0
	// number + step >  limit ?	number = number + step - limit = number - (limit - step)
	if c.number >= (c.settings.Limit - c.settings.Step) {
		// переход счетчика через границу
		// делаем за две операции, чтобы не выйти за разрядность
		rest := c.settings.Limit - c.settings.Step
		c.number -= rest
	} else {
		c.number += c.settings.Step
	}
	return &c.number, nil
}

// checkSettings used internally for checking if settings correct
// ie Step < Limit and Step > 0
func checkSettings(s *setup.Settings) (e error) {
	if s.Step >= s.Limit {
		e = errors.New("Step must be less than limit")
	} else if s.Step <= 0 {
		e = errors.New("Step must be positive")
	}
	return
}
