package filesystem

import (
	"path/filepath"
	"sort"

	"github.com/jbystronski/godirscan/pkg/common"
)

type Entries []*FsFiletype

func (e *Entries) AllItems() []*common.StoreItem {
	items := common.Entries{}

	for _, v := range e.All() {

		item := common.StoreItem(*v)
		items = append(items, &item)

	}

	return items
}

func (e *Entries) FindItem(index int) *common.StoreItem {
	item := e.Find(index)

	sItem := (*item).(common.StoreItem)

	return &sItem
}

func (e *Entries) All() []*FsFiletype {
	return *e
}

func (e Entries) Last() int {
	return len(e.All()) - 1
}

func (e Entries) Len() int {
	return len(e.All())
}

func (e *Entries) Swap(i, j int) {
	(*e)[i], (*e)[j] = (*e)[j], (*e)[i]
}

func (e *Entries) Less(i, j int) bool {
	return (*e.Find(i)).Size() < (*e.Find(j)).Size()
}

func (e *Entries) Find(index int) *FsFiletype {
	if index > len(*e)-1 {
		return nil
	}

	return (*e)[index]
}

func (e *Entries) FindByPath(path string) *FsFiletype {
	for _, en := range e.All() {
		if (*en).FullPath() == path {
			return en
		}
	}

	return nil
}

func (e *Entries) Insert(entry FsFiletype) {
	*e = append(e.All(), &entry)
}

func (e *Entries) SortByName() {
	sort.Slice(e.All(), func(i, j int) bool {
		left := (*e)[i]
		right := (*e)[j]
		_, isLeftDir := (*left).(*FsDirectory)
		_, isRightDir := (*right).(*FsDirectory)

		if isLeftDir && !isRightDir {
			return true
		} else if !isLeftDir && isRightDir {
			return false
		} else {
			return (*left).Name() < (*right).Name()
		}
	})
}

func (e *Entries) sortByType() {
	sort.Slice(e.All(), func(i, j int) bool {
		left := (*e)[i]
		right := (*e)[j]
		_, isLeftDir := (*left).(*FsDirectory)
		_, isRightDir := (*right).(*FsDirectory)

		if isLeftDir && !isRightDir {
			return true
		} else if !isLeftDir && isRightDir {
			return false
		} else {
			return filepath.Ext((*left).Name()) < filepath.Ext((*right).Name())
		}
	})
}
