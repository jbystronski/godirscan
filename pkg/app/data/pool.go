package data

import "sync"

type Pool struct {
	pool sync.Pool
}

func NewDataPool() *Pool {
	entity := &Pool{}

	entity.pool = sync.Pool{
		New: func() interface{} {
			return &FsEntry{}
		},
	}

	return entity
}

func (p *Pool) Get() *FsEntry {
	entry := p.pool.Get().(*FsEntry)
	return entry
}

func (p *Pool) Put(e *FsEntry) {
	p.pool.Put(e)
}
