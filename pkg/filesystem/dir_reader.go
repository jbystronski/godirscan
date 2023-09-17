package filesystem

import (
	"io/fs"
	"os"
)

type DirReader struct{}

func (d *DirReader) Read(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (d *DirReader) ReadIgnorePermission(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsPermission(err) {
			return []fs.DirEntry{}, nil
		} else {
			return nil, err
		}
	}

	return entries, nil
}
