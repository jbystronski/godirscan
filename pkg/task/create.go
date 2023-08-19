package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/terminal"
)

func CreateFsFile(path string, offset int) (ok bool, err error) {
	name, inputErr := WaitInput("Create a new file: ", "", terminal.PromptLine, offset)
	if inputErr != nil {
		err = inputErr
		return
	}

	if name == "" {
		return
	}

	name = strings.TrimSpace(name)
	if strings.ContainsAny(name, string(os.PathSeparator)) {
		err = fmt.Errorf("%s: %v", "filename contains path separator, aborting", os.PathSeparator)

		return
	}
	newFilePath := filepath.Join(path, name)

	_, statErr := os.Stat(newFilePath)

	if statErr != nil && errors.Is(statErr, os.ErrNotExist) {

		_, createErr := os.Create(newFilePath)
		if createErr != nil {
			return false, createErr
		}

	} else {

		answ, inputErr := WaitInput(fmt.Sprintf("%s (%s) %s", "File", name, "already exists, do you wish to override it?"), "n", terminal.PromptLine, offset)

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

func CreateFsDirectory(path string, offset int) (ok bool, err error) {
	dir, err := WaitInput("Create directory: ", "", terminal.PromptLine, offset)
	if err != nil {
		return ok, err
	}

	if dir == "" {
		return
	}

	dir = strings.TrimSpace(dir)
	if strings.ContainsAny(dir, string(os.PathSeparator)) {
		err = fmt.Errorf("%s \"%v\"", "Folder name cannot contain", string(os.PathSeparator))

		return
	}
	newDirPath := filepath.Join(path, dir)
	err = os.Mkdir(newDirPath, 0o777)
	if err != nil {
		return
	}
	ok = true
	return
}
