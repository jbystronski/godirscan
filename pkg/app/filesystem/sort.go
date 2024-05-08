package filesystem

import (
	"sort"
)

func (c *FsController) sort() {
	if c.defaultSort == nil {
		c.defaultSort = c.sortByType
	}

	c.defaultSort()
}

func (f *FsController) sortByName() {
	f.data.SortByName()
	f.setDefaultSort(f.sortByName)
}

func (f *FsController) sortByType() {
	f.data.SortByType()
	f.setDefaultSort(f.sortByType)
}

func (f *FsController) sortBySizeAsc() {
	sort.Sort(f.data)
	f.setDefaultSort(f.sortBySizeAsc)
}

func (f *FsController) sortBySizeDesc() {
	sort.Sort(sort.Reverse(f.data))
	f.setDefaultSort(f.sortBySizeDesc)
}

func (f *FsController) setDefaultSort(fn func()) {
	f.defaultSort = fn
	f.fullRender()
}
