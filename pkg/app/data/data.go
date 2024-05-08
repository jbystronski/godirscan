package data

import (
	"path/filepath"
	"sort"
)

type FsData struct {
	data []*FsEntry
}

func NewFsData() *FsData {
	entity := &FsData{}
	entity.data = []*FsEntry{}

	return entity
}

func (e *FsData) All() []*FsEntry {
	return e.data
}

func (e FsData) Size() int {
	var size int

	for _, v := range e.All() {
		size += v.Size()
	}

	return size
}

func (e *FsData) Reset() {
	e.data = []*FsEntry{}
}

func (e FsData) Last() int {
	return len(e.All()) - 1
}

func (e FsData) Len() int {
	return len(e.All())
}

func (e FsData) Swap(i, j int) {
	e.data[i], e.data[j] = e.data[j], e.data[i]
}

func (e *FsData) Less(i, j int) bool {
	return e.All()[i].Size() < e.All()[j].Size()
}

func (e *FsData) Find(index int) (*FsEntry, bool) {
	if index > e.Len()-1 {
		return nil, false
	}

	return e.All()[index], true
}

func (e FsData) FindByPath(path string) (*FsEntry, bool) {
	for _, en := range e.All() {
		if en.FullPath() == path {
			return en, true
		}
	}

	return nil, false
}

func (e *FsData) Insert(entry *FsEntry) {
	e.data = append(e.data, entry)
}

func (e FsData) SortByName() {
	sort.Slice(e.All(), func(i, j int) bool {
		left := e.data[i]
		right := e.data[j]
		isLeftDir := left.FsType() == DirDatatype
		isRightDir := right.FsType() == DirDatatype

		if isLeftDir && !isRightDir {
			return true
		} else if !isLeftDir && isRightDir {
			return false
		} else {
			return left.Name() < right.Name()
		}
	})
}

func (e *FsData) SortByType() {
	sort.Slice(e.All(), func(i, j int) bool {
		left := e.data[i]
		right := e.data[j]
		isLeftDir := left.FsType() == DirDatatype
		isRightDir := right.FsType() == DirDatatype

		if isLeftDir && !isRightDir {
			return true
		} else if !isLeftDir && isRightDir {
			return false
		} else {
			return filepath.Ext(left.Name()) < filepath.Ext(right.Name())
		}
	})
}
