package filesystem

import "github.com/jbystronski/godirscan/pkg/common"

type DataAccessor interface {
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
}

type DataPool interface {
	Get() FsEntry
	Put(FsEntry)
}

type FsDataAccessor interface {
	Self() []*FsEntry
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
}

type PathAccessor interface {
	Path() string
	FullPath() string
}

type PathMutator interface {
	SetPath(string)
}

type FsFiletype interface {
	PathAccessor
	// Path() string
	// FullPath() string
	// // PathMutator
	common.SizeAccessor
	common.SizeMutator
	common.NameAccessor
	common.NameMutator

	String() string

	SetFsType(FsEntity)
	FsType() FsEntity
}

type FsStoreAccessor interface {
	common.NameAccessor
	common.NameMutator

	Data() FsDataAccessor
	SetData(string) error
	ResetData()

	GetParentDirectory() (string, bool)
	common.SizeAccessor
}
