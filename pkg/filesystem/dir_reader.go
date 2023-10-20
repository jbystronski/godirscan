package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
)

type DirReader struct{}

func (d *DirReader) Read(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (d *DirReader) GetParentDirectory(currentDir string) (string, bool) {
	if common.GetRootDirectory() == currentDir {
		return "", false
	}

	parent, _ := filepath.Split(currentDir)
	parent = strings.TrimSuffix(parent, string(filepath.Separator))

	if parent == "" {
		parent = common.GetRootDirectory()
	}

	return parent, true
}
