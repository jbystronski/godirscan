package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/utils"
)

func CreateFsFile(path string) (ok bool) {
	name := WaitInput("Create a new file: ", "")

	name = strings.TrimSpace(name)
	if strings.ContainsAny(name, string(os.PathSeparator)) {
		utils.ShowErrAndContinue(fmt.Errorf("%s: %v", "filename contains path separator, aborting", os.PathSeparator))

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

		ok = true

	} else {

		name = WaitInput(fmt.Sprintf("%s (%s) %s", "File", name, "already exists, do you wish to override it?"), "n")

		if name == "y" || name == strings.ToLower("YES") {
			os.Truncate(newFilePath, 0)

			ok = true

		}

	}

	return
}

func CreateFsDirectory(path string) (ok bool) {
	dir := WaitInput("Create directory: ", "")

	dir = strings.TrimSpace(dir)
	if strings.ContainsAny(dir, string(os.PathSeparator)) {
		utils.ShowErrAndContinue(fmt.Errorf("%s \"%v\"", "Folder name cannot contain", string(os.PathSeparator)))

		return
	}
	newDirPath := filepath.Join(path, dir)
	err := os.Mkdir(newDirPath, 0o777)
	if err != nil {
		utils.ShowErrAndContinue(err)
		return
	}

	ok = true

	return
}
