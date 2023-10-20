package common

import (
	"time"
)

type SizeAccessor interface {
	Size() int
}

type SizeMutator interface {
	SetSize(int)
}

type NameAccessor interface {
	Name() string
}

type NameMutator interface {
	SetName(string)
}

type Convert interface {
	Convert()
}

type Cancelable interface {
	Create()
	Cancel()
	IsCancelled() <-chan struct{}
}

type Tickable interface {
	Start(time.Duration)
	Stop()
	Interval() time.Duration
	SetInterval(time.Duration)
	Tick() <-chan time.Time
	IsInitialized() bool
}

type MapAccessor[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Unset(K)
	Clear()
	Exists(K) bool
	Len() int
	Self() map[K]V
	Copy() map[K]V
}
