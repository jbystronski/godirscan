package utils

import (
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
)

func GetParentDirectory(dir string) (string, bool) {
	if common.GetRootDirectory() == dir {
		return "", false
	}

	parent, _ := filepath.Split(dir)
	parent = strings.TrimSuffix(parent, string(filepath.Separator))

	if parent == "" {
		parent = common.GetRootDirectory()
	}

	return parent, true
}
