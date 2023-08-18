package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/utils"
)

func Rename(currentName, pathToFile string) (bool, error) {
	answ, err := WaitInput("rename", currentName)
	if err != nil {
		return false, err
	}

	if strings.Contains(answ, string(os.PathSeparator)) {
		utils.ShowErrAndContinue(fmt.Errorf("path separator can't be used inside name"))
		return false, nil

	}

	if answ == "" || answ == currentName {
		return false, nil
	}

	err = os.Rename(filepath.Join(pathToFile, currentName), filepath.Join(pathToFile, answ))

	if err != nil {
		return false, err
	}

	return true, nil
}
