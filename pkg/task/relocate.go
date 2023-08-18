package task

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/jbystronski/godirscan/pkg/config"
	"github.com/jbystronski/godirscan/pkg/navigator"
	"github.com/jbystronski/godirscan/pkg/terminal"
	"github.com/jbystronski/godirscan/pkg/utils"
)

var (
	stopProgress    = make(chan struct{}, 1)
	currentFileName = make(chan string)
	wg              sync.WaitGroup
	sem             = make(chan struct{}, config.Cfg.MaxWorkers)
	skipAll         bool
)

func writeFile(srcPath, targetPath string) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, srcFile)
	if err != nil {
		panic(err)
	}

	info, err := os.Stat(srcPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)

	}

	err = os.Chmod(targetPath, info.Mode())
	if err != nil {
		panic(err)
	}
}

func tryCreateSymlink(srcPath, targetPath string) bool {
	fileInfo, err := os.Lstat(srcPath)
	if err != nil {
		panic(err)
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		symlinkTarget, linkErr := os.Readlink(srcPath)

		if linkErr != nil {
			panic(linkErr)
		}

		symlinkErr := os.Symlink(symlinkTarget, targetPath)
		if symlinkErr != nil {
			panic(symlinkErr)
		}

		return true

	}
	return false
}

func move(srcPath, targetPath string, deleteCopied bool) {
	srcInfo, statErr := os.Stat(srcPath)
	if statErr != nil {

		wg.Add(1)

		go func(statErr error, srcPath, targetPath string) {
			defer func() {
				wg.Done()
			}()

			ok := tryCreateSymlink(srcPath, targetPath)
			<-sem

			if !ok {
				panic(statErr)
			}
		}(statErr, srcPath, targetPath)
		sem <- struct{}{}

	}

	if statErr == nil {
		if srcInfo.IsDir() {

			createErr := os.Mkdir(targetPath, srcInfo.Mode())
			if createErr != nil {
				panic(createErr)
			}
			dc, readErr := os.ReadDir(srcPath)
			if readErr != nil {
				panic(readErr)
			}

			for _, entry := range dc {
				move(filepath.Join(srcPath, entry.Name()), filepath.Join(targetPath, entry.Name()), deleteCopied)
			}

			if deleteCopied {
				removeErr := os.RemoveAll(srcPath)
				if removeErr != nil {
					panic(removeErr)
				}
			}
		} else {
			wg.Add(1)

			go func(srcPath, targetPath string) {
				defer func(path string) {
					wg.Done()
					currentFileName <- fmt.Sprint("Finished writing ", path, " ")
				}(targetPath)
				currentFileName <- fmt.Sprint("Writing ", targetPath, " ")
				writeFile(srcPath, targetPath)
				<-sem
			}(srcPath, targetPath)
			sem <- struct{}{}
		}
	}
}

func addSuffixIfExists(targetPath string) string {
	_, err := os.Stat(targetPath)
	if err != nil {
		return targetPath
	}

	for num := 1; ; num++ {
		targetCopy := targetPath + " COPY " + strconv.Itoa(num)
		_, err := os.Stat(targetCopy)
		if err != nil {
			return targetCopy
		}
	}
}

func Relocate(prompt string, deleteCopied bool, selected *navigator.Selected, currentPath string) (bool, error) {
	answ, err := WaitInput(fmt.Sprintf("%s %s", prompt, "selected into the current directory? :"), "y")
	if err != nil {
		return false, err
	}

	if answ != "y" {
		return false, nil
	}

	go func() {
		terminal.PrintProgress(stopProgress, currentFileName)
	}()

	for entry := range selected.GetAll() {

		if *entry.Path == currentPath {
			stopProgress <- struct{}{}
			utils.ShowErrAndContinue(errors.New("copying / moving within same directory is not permitted"))

			return false, nil
		}

		if strings.HasPrefix(currentPath, entry.FullPath()) {
			stopProgress <- struct{}{}
			utils.ShowErrAndContinue(errors.New("cannot move / copy a folder into itself"))

			return false, nil
		}

		wg.Add(1)
		srcPath := filepath.Join(*entry.Path, entry.Name)

		targetPath := addSuffixIfExists(filepath.Join(currentPath, entry.Name))
		go func(srcPath, targetPath string) {
			defer wg.Done()

			terminal.ClearLine()
			terminal.CarriageReturn()

			move(srcPath, targetPath, deleteCopied)
		}(srcPath, targetPath)

	}
	wg.Wait()

	terminal.ClearLine()
	terminal.CarriageReturn()
	stopProgress <- struct{}{}
	return true, nil

	// TODO: after copying sort according to current sorting alghoritm
}
