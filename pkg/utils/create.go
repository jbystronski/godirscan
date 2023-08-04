package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/entry"
)

func createFsFile(path string, allEntries *[]*entry.Entry, errHandle func(error)) (ok bool) {
	waitUserInput("Create a new file: ", "", func(name string) {
		name = strings.TrimSpace(name)
		if strings.ContainsAny(name, string(os.PathSeparator)) {
			errHandle(fmt.Errorf("%s: %v", "filename contains path separator, aborting", os.PathSeparator))

			return
		}
		newFilePath := filepath.Join(path, name)

		_, err := os.Stat(newFilePath)

		if err != nil && errors.Is(err, os.ErrNotExist) {

			_, err := os.Create(newFilePath)
			if err != nil {
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}

			info, err := os.Stat(newFilePath)

			*allEntries = append(*allEntries, &entry.Entry{Name: info.Name(), Size: int(info.Size()), IsDir: false, Path: &path})
			ok = true

		} else {
			waitUserInput(fmt.Sprintf("%s (%s) %s", "File", name, "already exists, do you wish to override it?"), "n", func(answ string) {
				if answ == "y" || answ == strings.ToLower("YES") {
					os.Truncate(newFilePath, 0)
					for _, en := range *allEntries {
						if en.Name == name {
							en.Size = 0
						}
					}
					ok = true

				}
			})
		}
	})
	return
}

func createFsDirectory(path string, allEntries *[]*entry.Entry, errHandle func(error)) (ok bool) {
	waitUserInput("Create directory: ", "", func(name string) {
		name = strings.TrimSpace(name)
		if strings.ContainsAny(name, string(os.PathSeparator)) {
			errHandle(fmt.Errorf("%s \"%v\"", "Folder name cannot contain", string(os.PathSeparator)))

			return
		}
		newDirPath := filepath.Join(path, name)
		err := os.Mkdir(newDirPath, 0o777)
		if err != nil {
			errHandle(err)
			return
		}

		info, err := os.Stat(newDirPath)
		if err != nil {
			errHandle(err)
			return
		}

		*allEntries = append(*allEntries, &entry.Entry{Name: name, IsDir: true, Size: int(info.Size()), Path: &path})
		ok = true

		// dir.Size += entry.Size
	})

	return
}
