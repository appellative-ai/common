package core

import (
	"sync"
)

// MapT - map type
type MapT[T, U any] struct {
	m *sync.Map
}

// NewSyncMap - create a new map
func NewSyncMap[T, U any]() *MapT[T, U] {
	c := new(MapT[T, U])
	c.m = new(sync.Map)
	return c
}

func (m *MapT[T, U]) Load(t T) (u U, ok bool) {
	v, ok1 := m.m.Load(t)
	if !ok1 {
		return u, ok1
	}
	if v1, ok2 := v.(U); ok2 {
		return v1, ok2
	}
	return u, false
}

func (m *MapT[T, U]) Store(t T, u U) {
	m.m.Store(t, u)
}
