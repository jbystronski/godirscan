package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func scanInputDirectory(defaultDir string, errHandle func(err error)) (rootDir string, entries []*entry.Entry) {
	waitUserInput("Scan directory: ", defaultDir, func(fPath string) {
		// if strings.HasPrefix(fPath, "~") {
		// 	currentUser, err := user.Current()
		// 	printDefaultErrorAndExit(err)

		// 	fPath = strings.Replace(fPath, "~", currentUser.HomeDir, 1)

		// }

		terminal.ResolveUserDirectory(&fPath)

		_, err := os.Stat(fPath)
		if err != nil {
			errHandle(err)

			return

		}

		// name := en.Name()
		rootPath := filepath.VolumeName("") + string(os.PathSeparator)

		if fPath == rootPath {
			fPath = rootPath
			// name = ""
		} else {
			fPath = strings.TrimSuffix(fPath, string(os.PathSeparator))

			// splitPath := strings.Split(fPath, string(os.PathSeparator))
			// splitPath = splitPath[:len(splitPath)-1]
			// fPath = string(os.PathSeparator) + filepath.Join(splitPath...)
		}

		// root = &Entry{
		// 	Name:  name,
		// 	Size:  int(en.Size()),
		// 	IsDir: en.IsDir(),

		// 	Path: fPath,
		// }

		rootDir, entries = scanDirectory(fPath)

		// root.Name = name
		// root.IsDir = en.IsDir()
		// root.Path = fPath
	})

	return
}

func scanDirectory(path string) (string, []*entry.Entry) {
	allEntries := []*entry.Entry{}

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("error during os.stat: ", err)
	}

	if !info.IsDir() {

		fmt.Println("error occcured")
		return "", nil
	}

	dc, err := os.ReadDir(path)
	if err != nil {

		fmt.Println("error occured")
		return "", nil
	}

	for _, en := range dc {

		info, err := os.Lstat(filepath.Join(path, en.Name()))
		if err != nil {

			fmt.Println("error occurred")
			continue
		}

		var size int
		var isDir bool

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			size = 0
			isDir = false
		} else {
			size = int(info.Size())
			isDir = info.IsDir()
		}

		newEntry := &entry.Entry{
			Name:  en.Name(),
			Size:  size,
			IsDir: isDir,
			Path:  &path,
		}

		allEntries = append(allEntries, newEntry)
	}

	return path, allEntries
}

// func scanDirectorySize(allEntries *[]*Entry, totalSize *int) {
// 	maxWorkers := Cfg.MaxWorkers

// 	if maxWorkers == 0 {
// 		maxWorkers = defaultConfig.MaxWorkers
// 	}

// 	var scanWg sync.WaitGroup
// 	var calculateSize func(string, *Entry, chan struct{}, *sync.WaitGroup)
// 	// ch := make(chan int)
// 	workerChan := make(chan struct{}, maxWorkers)

// 	calculateSize = func(path string, topDir *Entry, workerChan chan struct{}, wg *sync.WaitGroup) {
// 		dc, err := os.ReadDir(path)
// 		if err != nil {
// 			// fmt.Println("Error reading directory: ", err)
// 			// fmt.Println(len(workerChan))
// 			// printDefaultErrorAndExit(err)
// 			// return
// 		} else {
// 			for _, en := range dc {

// 				info, err := en.Info()
// 				if err != nil {
// 					// fmt.Println("Error during en.Info")

// 					// fmt.Println("Error traversing directory: ", err)

// 					// fmt.Println(len(workerChan))

// 					continue
// 				} else {
// 					topDir.Size += int(info.Size())

// 					if info.IsDir() {
// 						// calculateSize(filepath.Join(path, info.Name()), topDir, wg)

// 						if len(workerChan) >= maxWorkers {
// 							time.Sleep(time.Millisecond * 100)
// 						}

// 						workerChan <- struct{}{}
// 						// fmt.Println(len(workerChan))

// 						wg.Add(1)

// 						subPath := filepath.Join(path, en.Name())

// 						go func(subPath string, topDir *Entry, workerChan chan struct{}, wg *sync.WaitGroup) {
// 							defer func() {
// 								<-workerChan
// 								// fmt.Println(len(workerChan))

// 								wg.Done()
// 							}()

// 							calculateSize(subPath, topDir, workerChan, wg)
// 						}(subPath, topDir, workerChan, wg)
// 					}
// 				}

// 			}
// 		}
// 	}

// 	for _, en := range *allEntries {
// 		if en.IsDir {
// 			calculateSize(filepath.Join(*en.Path, en.Name), en, workerChan, &scanWg)
// 		}
// 	}

// 	scanWg.Wait()

// 	close(workerChan)

// 	for _, v := range *allEntries {
// 		*totalSize += v.Size
// 	}

// 	// doneChan <- struct{}{}
// }

var virtualFsMap = map[string]struct{}{
	"/proc": {},
	"/dev":  {},
	"/sys":  {},
}

func isVirtualFs(name string) bool {
	for virtualFsEntry := range virtualFsMap {
		if strings.HasPrefix(name, virtualFsEntry) {
			return true
		}
	}

	return false
}

func scanDirectorySize(allEntries *[]*entry.Entry, totalSize *int) {
	*totalSize = 0
	var calculateSize func(string, *entry.Entry)
	// ch := make(chan int)

	var wg sync.WaitGroup

	calculateSize = func(path string, topDir *entry.Entry) {
		dc, err := os.ReadDir(path)
		if err != nil {
			// fmt.Println("Error reading directory: ", err)
			// fmt.Println(len(workerChan))
			// printDefaultErrorAndExit(err)
			// return
		} else {
			for _, en := range dc {

				info, err := en.Info()
				if err != nil {
					// fmt.Println("Error during en.Info")

					// fmt.Println("Error traversing directory: ", err)

					// fmt.Println(len(workerChan))

					continue
				} else {
					topDir.Size += int(info.Size())

					if info.IsDir() {
						// calculateSize(filepath.Join(path, info.Name()), topDir, wg)

						// fmt.Println(len(workerChan))

						subPath := filepath.Join(path, en.Name())

						calculateSize(subPath, topDir)

					}
				}

			}
		}
	}

	for _, en := range *allEntries {
		if en.IsDir {

			path := filepath.Join(*en.Path, en.Name)

			if isVirtualFs(path) {
				en.Size = 0
				continue
			}

			wg.Add(1)
			go func(en *entry.Entry) {
				defer wg.Done()
				calculateSize(filepath.Join(*en.Path, en.Name), en)
			}(en)
		}
	}

	wg.Wait()

	for _, v := range *allEntries {
		*totalSize += v.Size
	}

	// doneChan <- struct{}{}
}
