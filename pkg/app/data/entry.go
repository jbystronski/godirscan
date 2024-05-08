package data

import (
	"fmt"
	"os"
	"path/filepath"
)

type FsEntity uint8

const (
	FileDatatype FsEntity = iota
	DirDatatype
	SymlinkDatatype
	SearchResultDatatype
)

type FsEntry struct {
	name   string
	size   int
	fsType FsEntity
}

func (e *FsEntry) SetName(n string) {
	e.name = n
}

func (e *FsEntry) Name() string {
	return filepath.Base(e.name)
}

func (e FsEntry) Size() int {
	return e.size
}

func (e *FsEntry) SetSize(size int) {
	e.size = size
}

func (e *FsEntry) FullPath() string {
	return fmt.Sprint(filepath.Join(e.Path(), e.Name()))
}

func (e *FsEntry) Path() string {
	return filepath.Dir(e.name)
}

func (e *FsEntry) Rename(newName string) error {
	err := os.Rename(e.FullPath(), filepath.Join(e.Path(), newName))
	if err != nil {
		return err
	}

	e.SetName(filepath.Join(e.Path(), newName))

	return nil
}

func (e *FsEntry) FsType() FsEntity {
	return e.fsType
}

func (e *FsEntry) SetFsType(t FsEntity) {
	e.fsType = t
}

func (e FsEntry) String() string {
	switch e.FsType() {

	case SearchResultDatatype:
		return e.FullPath()
	case SymlinkDatatype:
		return fmt.Sprint(e.Name(), " ", "@")
	default:
		return e.Name()

	}
}
