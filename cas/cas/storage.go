package cas

import (
	"errors"
	"sync"
)

var errMismatch = errors.New("invalid ver")

type Storage struct {
	value Value
	mu    sync.Mutex
}

type Value struct {
	val     int
	version int
}

func NewStorage(val int) *Storage {
	return &Storage{
		value: Value{
			val:     val,
			version: 0,
		},
	}
}

func (a *Storage) Get() Value {
	return a.value
}

func (a *Storage) Set(newValue Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.value.version != newValue.version {
		return errMismatch
	}
	newValue.version++
	a.value = newValue
	return nil
}
