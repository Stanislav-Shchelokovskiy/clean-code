package race

import "sync"

type Storage struct {
	value Value
	mu    sync.Mutex
}

type Value struct {
	val int
}

func NewStorage(val int) *Storage {
	return &Storage{
		value: Value{val: val},
	}
}

func (a *Storage) Get() Value {
	return a.value
}

func (a *Storage) Set(newValue Value) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.value = newValue
}
