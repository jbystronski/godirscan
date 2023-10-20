package utils

import (
	"io/fs"
	"os"
)

func IsSymlink(info fs.FileInfo) bool {
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}

	return false
}
