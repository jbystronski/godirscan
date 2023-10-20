package utils

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func Search(path, pattern string, errChan chan<- error, entryChan chan<- struct {
	Path string
	Info fs.FileInfo
}, test func(fs.FileInfo) bool,
) {
	if dirEntries, err := ReadIgnorePermission(path); err != nil {

		errChan <- err
		return

	} else {
		for _, en := range dirEntries {
			info, err := en.Info()
			if err != nil {

				errChan <- err

				return
			}

			if info.IsDir() {
				Search(filepath.Join(path, en.Name()), pattern, errChan, entryChan, test)
			} else {
				if strings.Contains(info.Name(), pattern) {
					if ok := test(info); ok {
						entryChan <- struct {
							Path string
							Info fs.FileInfo
						}{filepath.Join(path, info.Name()), info}
					}
				}
			}

		}
	}
}
