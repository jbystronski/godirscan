package common

import "sync"

var rw sync.RWMutex

type GenericMap[K comparable, V any] map[K]V

func (m *GenericMap[K, V]) Clear() {
	*m = make(map[K]V)
}

func (m GenericMap[K, V]) Get(key K) (V, bool) {
	rw.RLock()
	defer rw.RUnlock()
	if value, exist := (m)[key]; exist {
		return value, true
	}

	var none V
	return none, false
}

func (m GenericMap[K, V]) Exists(key K) bool {
	if _, ok := m.Get(key); ok {
		return true
	}

	return false
}

func (m *GenericMap[K, V]) Set(key K, value V) {
	rw.Lock()
	defer rw.Unlock()

	(*m)[key] = value
}

func (m *GenericMap[K, V]) Unset(key K) {
	delete(*m, key)
}

func (m GenericMap[K, V]) Len() int {
	return len(m)
}

func NewGenericMap[K comparable, V any]() GenericMap[K, V] {
	m := make(GenericMap[K, V])
	return m
}

func (m *GenericMap[K, V]) Self() map[K]V {
	return *m
}

func (m *GenericMap[K, V]) Copy() map[K]V {
	copy := make(GenericMap[K, V])

	for k, v := range m.Self() {
		copy[k] = v
	}

	return copy
}
