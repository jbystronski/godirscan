package filesystem

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
)

type Finder struct {
	common.Cancelable
	common.Tickable
	DirReader
}

func (f *Finder) find(path string, pattern string, errChan chan<- error, doneChan chan<- struct{}, addEntry func(FsFiletype), test func(fs.FileInfo) bool) {
	f.Cancel()

	f.Tickable.Start(f.Interval())
	f.Cancelable.Create()

	go func(path string) {
		defer f.Tickable.Stop()
		f.runSearch(filepath.Join(path), pattern, errChan, addEntry, test)

		doneChan <- struct{}{}
	}(path)
}

func (f *Finder) runSearch(path, pattern string, errChan chan<- error, addEntry func(FsFiletype), test func(fs.FileInfo) bool) {
	go func() {
		for {
			select {
			case <-f.IsCancelled():

				return
			}
		}
	}()

	if dirEntries, err := f.ReadIgnorePermission(path); err != nil {

		errChan <- err
		return

	} else {
		var newEntry FsFiletype

		for _, en := range dirEntries {
			info, err := en.Info()
			if err != nil {
				common.Log("finder en.Info err ", err)
				errChan <- err

				return
			}

			if info.IsDir() {
				f.runSearch(filepath.Join(path, en.Name()), pattern, errChan, addEntry, test)
			} else {
				if strings.Contains(info.Name(), pattern) {
					if ok := test(info); ok {
						newEntry = &FsFound{}

						newEntry.SetName(info.Name())
						newEntry.SetSize(int(info.Size()))
						newEntry.SetPath(path)
						addEntry(newEntry)

					}
				}
			}

		}
	}
}
