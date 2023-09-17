package filesystem

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jbystronski/godirscan/pkg/common"
)

var virtualFsMap = map[string]struct{}{
	"/proc": {},
	"/dev":  {},
	"/sys":  {},
}

type Scanner struct {
	common.Cancelable
	common.Tickable
	DirReader
}

func (s *Scanner) isVirtualFilesystem(dirName string) bool {
	for virtualFsEntry := range virtualFsMap {
		if strings.HasPrefix(dirName, virtualFsEntry) {
			return true
		}
	}

	return false
}

func (s *Scanner) GetPathInfo(path string) (fs.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *Scanner) ResolveUserDirectory(path *string) {
	switch runtime.GOOS {
	case "darwin":
		{
			break
		}
	case "windows":
		{
			break
		}
	default:
		{
			if strings.HasPrefix(*path, "~") {
				currentUser, err := user.Current()
				if err != nil {
					panic(err)
				}

				*path = strings.Replace(*path, "~", currentUser.HomeDir, 1)

			}
		}
	}
}

func (s *Scanner) scan(path string) (*Entries, error) {
	if dc, err := s.ReadIgnorePermission(path); err != nil {
		return nil, err
	} else {

		entries := Entries{}

		var newEntry FsFiletype
		for _, en := range dc {

			info, err := os.Lstat(filepath.Join(path, en.Name()))
			if err != nil {
				continue
			}

			switch true {

			case s.IsSymlink(info):
				newEntry = &FsSymlink{}
				newEntry.SetSize(0)
			case info.IsDir():
				newEntry = &FsDirectory{}
				newEntry.SetSize(int(info.Size()))
			default:
				newEntry = &FsFile{}
				newEntry.SetSize(int(info.Size()))
			}

			newEntry.SetName(info.Name())

			newEntry.SetPath(path)

			entries.Insert(newEntry)

		}

		return &entries, nil

	}

	// info, err := os.Stat(path)
	// if err != nil {
	// 	return nil, err
	// }

	// if !info.IsDir() {
	// 	return nil, errors.New("Path is not a directory")
	// }

	// dc, err := os.ReadDir(path)
	// if err != nil {
	// 	return nil, err
	// }
}

func (s *Scanner) IsSymlink(info fs.FileInfo) bool {
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}

	return false
}

func (s *Scanner) scanDirectorySize(file *FsFiletype, path string, errChan chan<- error, wg *sync.WaitGroup, sem chan struct{}) {
	go func() {
		for {
			select {
			case <-s.IsCancelled():

				return
			}
		}
	}()

	if dc, err := s.ReadIgnorePermission(path); err != nil {
		errChan <- err
		return
	} else {
		for _, en := range dc {

			info, err := en.Info()
			if err != nil {
				errChan <- err
				return
			}

			// if !info.IsDir() {
			// 	size := (*file).Size() + int(info.Size())
			// }

			if !info.IsDir() {

				size := (*file).Size() + int(info.Size())
				(*file).SetSize(size)
			}

			//	size := int(info.Size())

			if info.IsDir() {
				wg.Add(1)

				go func(path string, info fs.FileInfo) {
					defer func() {
						wg.Done()

						<-sem
					}()

					s.scanDirectorySize(file, filepath.Join(path, info.Name()), errChan, wg, sem)
				}(path, info)

				sem <- struct{}{}
			}

		}
	}
}

//   // Create a worker function.
//   func worker(path string, info fs.FileInfo, file *FsFiletype, errChan chan<- error, wg *sync.WaitGroup, sem chan struct{}) {
// 		defer func() {
// 				wg.Done()
// 				<-sem
// 		}()
// 		s.scanDirectorySize(file, filepath.Join(path, info.Name()), errChan, wg, sem)
// }

// // Create a worker pool.
// numWorkers := 10 // Adjust this based on your system's capabilities.
// for i := 0; i < numWorkers; i++ {
// 		go func() {
// 				for {
// 						select {
// 						case <-s.IsCancelled():
// 								return
// 						case task, ok := <-workerQueue:
// 								if !ok {
// 										return
// 								}
// 								worker(task.path, task.info, task.file, task.errChan, task.wg, task.sem)
// 						}
// 				}
// 		}()
// }

// // Enqueue tasks for scanning directories.
// for _, en := range dc {
// 		if info.IsDir() {
// 				wg.Add(1)
// 				sem <- struct{}{}
// 				workerQueue <- WorkerTask{
// 						path:    path,
// 						info:    info,
// 						file:    file,
// 						errChan: errChan,
// 						wg:      wg,
// 						sem:     sem,
// 				}
// 		}
// }

// // Close workerQueue and wait for workers to finish.
// close(workerQueue)
// wg.Wait()

func (s *Scanner) scanDataSize(entries []*FsFiletype) error {
	var err error
	s.Cancel()
	s.Tickable.Start(s.Interval())
	s.Cancelable.Create()
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	doneChan := make(chan struct{})
	sem := make(chan struct{}, 550)

	defer func() {
		s.Tickable.Stop()
	}()

	go func() {
		for {
			select {
			case <-s.IsCancelled():

				return
			case err = <-errChan:
				s.Cancel()
			case <-doneChan:

				return

			}
		}
	}()

	for _, en := range entries {
		switch (*en).(type) {
		//	case *FsSymlink:
		//		(*en).SetSize(0)

		//	case *FsFile:
		// *totalSize += (*en).GetSize()
		case *FsDirectory:
			wg.Add(1)
			go func(en *FsFiletype) {
				defer func() {
					wg.Done()
				}()

				s.scanDirectorySize(en, (*en).FullPath(), errChan, &wg, sem)
			}(en)
		}
	}
	wg.Wait()

	doneChan <- struct{}{}
	return err
}

// defer wg.Done()

// 	for _, en := range entries {
// 		switch (*en).(type) {
// 		case *FsDirectory:

// 			if s.isVirtualFilesystem((*en).FullPath()) {

// 				(*en).SetSize(0)
// 				continue
// 			}
