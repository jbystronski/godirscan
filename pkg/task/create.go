package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/utils"
)

func CreateFsFile(path string) (ok bool, err error) {
	name, inputErr := WaitInput("Create a new file: ", "")
	if inputErr != nil {
		err = inputErr
		return
	}

	if name == "" {
		return
	}

	name = strings.TrimSpace(name)
	if strings.ContainsAny(name, string(os.PathSeparator)) {
		utils.ShowErrAndContinue(fmt.Errorf("%s: %v", "filename contains path separator, aborting", os.PathSeparator))

		return
	}
	newFilePath := filepath.Join(path, name)

	_, statErr := os.Stat(newFilePath)

	if statErr != nil && errors.Is(statErr, os.ErrNotExist) {

		_, createErr := os.Create(newFilePath)
		if createErr != nil {
			utils.ShowErrAndContinue(err)
			return
		}

	} else {

		answ, inputErr := WaitInput(fmt.Sprintf("%s (%s) %s", "File", name, "already exists, do you wish to override it?"), "n")

		if inputErr != nil {
			return ok, inputErr
		}

		if answ == "y" || answ == strings.ToLower("YES") {
			os.Truncate(newFilePath, 0)
		}

	}
	ok = true
	return ok, nil
}

func CreateFsDirectory(path string) (ok bool, err error) {
	dir, err := WaitInput("Create directory: ", "")
	if err != nil {
		return ok, err
	}

	dir = strings.TrimSpace(dir)
	if strings.ContainsAny(dir, string(os.PathSeparator)) {
		utils.ShowErrAndContinue(fmt.Errorf("%s \"%v\"", "Folder name cannot contain", string(os.PathSeparator)))

		return
	}
	newDirPath := filepath.Join(path, dir)
	err = os.Mkdir(newDirPath, 0o777)
	if err != nil {
		utils.ShowErrAndContinue(err)
		return
	}
	ok = true
	return
}
