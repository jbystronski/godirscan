package common

import (
	"time"
)

type StoreItem interface {
	NameAccessor
}

type DataAccessor interface {
	AllItems() []*StoreItem
	Len() int
	FindItem(int) *StoreItem
}

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

type StoreAccessor interface {
	NameAccessor
	NameMutator
	Items() DataAccessor
	SetItems(DataAccessor)
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

type Controller interface {
	ActionMap() map[ControllerAction]func()
	SetActionMap(map[ControllerAction]func())
	Map(ControllerAction)
	Active() bool
	SetActive(bool)
}
