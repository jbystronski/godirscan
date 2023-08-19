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
	"time"

	"github.com/jbystronski/godirscan/pkg/config"
	"github.com/jbystronski/godirscan/pkg/navigator"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

var (
	stopProgress    = make(chan struct{}, 1)
	currentFileName = make(chan string)
	wg              sync.WaitGroup
	sem             = make(chan struct{}, config.Cfg.MaxWorkers)
	skipAll         bool
)

func writeFile(srcPath, targetPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, srcFile)
	if err != nil {
		return err
	}

	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	err = os.Chmod(targetPath, info.Mode())
	if err != nil {
		return err
	}
	return nil
}

func tryCreateSymlink(srcPath, targetPath string) (bool, error) {
	fileInfo, err := os.Lstat(srcPath)
	if err != nil {
		return false, err
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		symlinkTarget, linkErr := os.Readlink(srcPath)

		if linkErr != nil {
			return false, linkErr
		}

		symlinkErr := os.Symlink(symlinkTarget, targetPath)
		if symlinkErr != nil {
			return false, symlinkErr
		}
		currentFileName <- fmt.Sprint("Created symlink ", targetPath, " ")
		return true, nil

	}
	return false, nil
}

func remove(srcPath string) error {
	removeErr := os.RemoveAll(srcPath)
	if removeErr != nil {
		return removeErr
	}
	return nil
}

func move(srcPath, targetPath string, deleteCopied bool) error {
	var err error
	var ok bool

	srcInfo, statErr := os.Stat(srcPath)
	if statErr != nil {

		wg.Add(1)

		go func(statErr error, srcPath, targetPath string) {
			defer func() {
				wg.Done()
			}()

			ok, err = tryCreateSymlink(srcPath, targetPath)
			if err != nil {
				return
			}
			<-sem

			if deleteCopied {
				err = remove(srcPath)
				if err != nil {
					return
				}
			}

			if !ok {
				err = statErr
			}
		}(statErr, srcPath, targetPath)
		sem <- struct{}{}

	}

	if statErr == nil {
		if srcInfo.IsDir() {

			err = os.Mkdir(targetPath, srcInfo.Mode())
			if err != nil {
				return err
			}
			dc, readErr := os.ReadDir(srcPath)
			if readErr != nil {
				err = readErr
				return err
			}

			for _, entry := range dc {
				err = move(filepath.Join(srcPath, entry.Name()), filepath.Join(targetPath, entry.Name()), deleteCopied)
				if err != nil {
					return err
				}

			}

			if deleteCopied {
				err = remove(srcPath)
				if err != nil {
					return err
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
				err = writeFile(srcPath, targetPath)

				if err != nil {
					return
				}

				if deleteCopied {
					err = remove(srcPath)
					if err != nil {
						return
					}
				}

				<-sem
			}(srcPath, targetPath)
			sem <- struct{}{}
		}
	}
	return err
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

func Relocate(prompt string, deleteCopied bool, selected *navigator.Selected, currentPath string, offset int) (bool, error) {
	answ, err := WaitInput(fmt.Sprintf("%s %s", prompt, "selected into the current directory? :"), "y", terminal.PromptLine, offset)
	if err != nil {
		panic(err)
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

			return false, errors.New("copying / moving within same directory is not permitted")
		}

		if strings.HasPrefix(currentPath, entry.FullPath()) {
			stopProgress <- struct{}{}

			return false, errors.New("cannot move / copy a folder into itself")
		}

		wg.Add(1)
		srcPath := filepath.Join(*entry.Path, entry.Name)

		targetPath := addSuffixIfExists(filepath.Join(currentPath, entry.Name))
		go func(srcPath, targetPath string) {
			defer wg.Done()

			terminal.ClearLine()
			terminal.CarriageReturn()

			err = move(srcPath, targetPath, deleteCopied)
			if err != nil {
				panic(err)
			}
		}(srcPath, targetPath)

	}
	wg.Wait()
	fmt.Println("FINISHED WAITGROUP")
	time.Sleep(time.Second * 2)
	terminal.ClearLine()
	terminal.CarriageReturn()
	stopProgress <- struct{}{}
	return true, nil

	// TODO: after copying sort according to current sorting alghoritm
}
