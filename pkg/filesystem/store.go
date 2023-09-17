package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
)

type FsStore struct {
	common.StoreAccessor

	FsDataAccessor
}

func (s *FsStore) Size() int {
	var size int

	for _, v := range s.Data().All() {
		size += (*v).Size()
	}

	return size
}

func (s *FsStore) Data() FsDataAccessor {
	return s.FsDataAccessor
}

func (s *FsStore) SetData(d FsDataAccessor) {
	s.FsDataAccessor = d
}

func NewStore(name string, accessor *Entries) *FsStore {
	store := &FsStore{
		FsDataAccessor: accessor,
	}

	store.SetName(name)

	return store
}

func (s *FsStore) getRootDirectory() string {
	currentDir, _ := os.Getwd()

	return filepath.VolumeName(currentDir) + string(filepath.Separator)
}

func (s *FsStore) ParentStoreName() (string, bool) {
	if s.getRootDirectory() == s.Name() {
		return "", false
	}

	parent, _ := filepath.Split(s.Name())
	parent = strings.TrimSuffix(parent, string(filepath.Separator))

	if parent == "" {
		parent = s.getRootDirectory()
	}

	return parent, true
}
