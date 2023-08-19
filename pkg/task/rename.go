package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/terminal"
)

func Rename(currentName, pathToFile string, offset int) (bool, error) {
	answ, err := WaitInput("rename", currentName, terminal.PromptLine, offset)
	if err != nil {
		return false, err
	}

	if strings.Contains(answ, string(os.PathSeparator)) {
		return false, fmt.Errorf("path separator can't be used inside name")
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
