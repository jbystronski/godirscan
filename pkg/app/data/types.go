package data

type Accessor interface {
	All() []*FsEntry
	Len() int
	Insert(*FsEntry)
	Less(i, j int) bool
	Swap(i, j int)
	Find(int) (*FsEntry, bool)
	SortByName()
	SortByType()
	FindByPath(string) (*FsEntry, bool)
	Reset()
	Size() int
	SetData([]*FsEntry)
}
