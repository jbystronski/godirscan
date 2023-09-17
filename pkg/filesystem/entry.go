package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
)

type FsEntry struct {
	common.Entry
	size int
	path string
}

func (e *FsEntry) SetName(n string) {
	e.Entry.SetName(n)
}

func (e *FsEntry) Size() int {
	return e.size
}

func (e *FsEntry) SetSize(size int) {
	e.size = size
}

func (e *FsEntry) FullPath() string {
	return fmt.Sprint(filepath.Join(e.path, e.Name()))
}

func (e *FsEntry) Path() string {
	return e.path
}

func (e *FsEntry) SetPath(p string) {
	e.path = p
}

func (e *FsEntry) Rename(newName string) (bool, error) {
	if newName == "" || newName == e.Name() {
		return false, nil
	}

	if strings.Contains(newName, string(os.PathSeparator)) {
		return false, fmt.Errorf("path separator can't be used inside name")
	}

	err := os.Rename(filepath.Join(e.Path(), e.Name()), filepath.Join(e.Path(), newName))
	if err != nil {
		return false, err
	}

	e.SetName(newName)

	return true, nil
}

func (e *FsEntry) printSize() string {
	return printSizeAsString(e.Size())
}
