package filesystem

import "github.com/jbystronski/godirscan/pkg/common"

type FsDataAccessor interface {
	All() []*FsFiletype
	Len() int
	Insert(FsFiletype)
	Less(i, j int) bool
	Swap(i, j int)
	Find(int) *FsFiletype
	SortByName()
	sortByType()
	FindByPath(string) *FsFiletype
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
	PathMutator
	common.SizeAccessor
	common.SizeMutator
	common.NameAccessor
	common.NameMutator
	Rename
	printSize() string
	String() string
}

type Rename interface {
	Rename(string) (bool, error)
}

type executable interface {
	execute() error
}

type Editable interface {
	Edit() error
}

type FsStoreAccessor interface {
	common.StoreAccessor
	Data() FsDataAccessor
	SetData(FsDataAccessor)

	ParentStoreName() (string, bool)
	common.SizeAccessor
}
