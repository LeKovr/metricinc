package counter

import (
	"errors"
	"sync"

	"lekovr/exam/counter/setup"
)

type Counter struct {
	number   int64
	settings setup.Settings
	mutex    sync.RWMutex // не указатель, т.к. всегда должен быть инициализирован
}

func NewCounter(s setup.Settings, number int64) (*Counter, error) {

	if err := checkSettings(s); err != nil {
		return nil, err
	}
	cnt := Counter{number: number, settings: s}
	return &cnt, nil
}

func (c *Counter) GetNumber() (int64, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.number, nil
}

func (c *Counter) IncrementNumber() (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	return c.number, nil
}

func (c *Counter) SetSettings(s setup.Settings) error {
	if err := checkSettings(s); err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.settings = s
	return nil
}

func (c *Counter) GetSettings() (setup.Settings, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.settings, nil
}

func checkSettings(s setup.Settings) (e error) {
	if s.Step >= s.Limit {
		e = errors.New("Step must be less than limit")
	}
	return
}