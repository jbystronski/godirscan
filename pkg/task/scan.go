package task

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jbystronski/godirscan/pkg/entry"
)

func resolveUserDirectory(fPath *string) {
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
			if strings.HasPrefix(*fPath, "~") {
				currentUser, err := user.Current()
				if err != nil {
					panic(err)
				}

				*fPath = strings.Replace(*fPath, "~", currentUser.HomeDir, 1)

			}
		}
	}
}

func ScanInputDirectory(defaultDir string) (rootDir string, entries []*entry.Entry) {
	WaitUserInput("Scan directory: ", defaultDir, func(fPath string) {
		resolveUserDirectory(&fPath)

		_, err := os.Stat(fPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)

		}

		rootPath := filepath.VolumeName("") + string(os.PathSeparator)

		if fPath == rootPath {
			fPath = rootPath
		} else {
			fPath = strings.TrimSuffix(fPath, string(os.PathSeparator))
		}

		rootDir, entries = ScanDirectory(fPath)
	})

	return
}

func ScanDirectory(path string) (string, []*entry.Entry) {
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

func ScanDirectorySize(allEntries []*entry.Entry, totalSize *int) {
	*totalSize = 0
	var calculateSize func(string, *entry.Entry)

	var wg sync.WaitGroup

	calculateSize = func(path string, topDir *entry.Entry) {
		dc, err := os.ReadDir(path)
		if err != nil {
		} else {
			for _, en := range dc {

				info, err := en.Info()
				if err != nil {
					continue
				} else {
					topDir.Size += int(info.Size())

					if info.IsDir() {

						subPath := filepath.Join(path, en.Name())

						calculateSize(subPath, topDir)

					}
				}

			}
		}
	}

	for _, en := range allEntries {
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

	for _, v := range allEntries {
		*totalSize += v.Size
	}
}
