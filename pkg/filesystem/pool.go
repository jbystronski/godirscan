package filesystem

import "sync"

type FsDataPool struct {
	pool sync.Pool
}

func NewDataPool() *FsDataPool {
	entity := &FsDataPool{}

	entity.pool = sync.Pool{
		New: func() interface{} {
			return &FsEntry{}
		},
	}

	return entity
}

func (p *FsDataPool) Get() *FsEntry {
	entry := p.pool.Get().(*FsEntry)
	return entry
}

func (p *FsDataPool) Put(e *FsEntry) {
	p.pool.Put(e)
}
