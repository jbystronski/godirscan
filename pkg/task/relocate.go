package task

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func writeFile(srcPath, srcName, targetPath string) {
	newFilepath := filepath.Join(targetPath, srcName)

	srcFile, err := os.Open(filepath.Join(srcPath, srcName))
	if err != nil {
		panic(err)
	}

	defer srcFile.Close()

	targetFile, err := os.Create(newFilepath)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, srcFile)
	if err != nil {
		panic(err)
	}

	info, err := os.Stat(filepath.Join(srcPath, srcName))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)

	}

	err = os.Chmod(newFilepath, info.Mode())
	if err != nil {
		panic(err)
	}
}

func move(srcPath, srcName, targetDir string, deleteCopied bool) {
	srcInfo, err := os.Stat(filepath.Join(srcPath, srcName))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if srcInfo.IsDir() {
		proceed := true

		err := os.Mkdir(filepath.Join(targetDir, srcName), srcInfo.Mode())
		if err != nil {
			if errors.Is(err, os.ErrExist) {

				answ := WaitInput(fmt.Sprintf("%s %s %s", "Folder ", srcName, " already exists, do you wish to merge them?"), "y")

				if answ == "y" || answ == strings.ToLower("YES") {
					proceed = true
				} else {
					proceed = false
				}

			}
		}

		if proceed {

			dc, err := os.ReadDir(filepath.Join(srcPath, srcName))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)

			}

			for _, entry := range dc {
				move(filepath.Join(srcPath, srcName), entry.Name(), filepath.Join(targetDir, srcName), deleteCopied)
			}

		}

		if deleteCopied {
			err := os.RemoveAll(filepath.Join(srcPath, srcName))
			if err != nil {
				panic(err)
			}
		}

	} else {

		_, err := os.Stat(filepath.Join(targetDir, srcName))

		if err != nil && errors.Is(err, os.ErrNotExist) {
			wg.Add(1)
			sem <- struct{}{}
			go func(srcPath, srcName, targetDir string) {
				defer func() {
					wg.Done()
					currentFileName <- fmt.Sprint("Written ", filepath.Join(targetDir, srcName), " ")
					<-sem
				}()

				writeFile(srcPath, srcName, targetDir)
			}(srcPath, srcName, targetDir)

		} else {

			answ := WaitInput(fmt.Sprintf("%s %s %s", "File", srcName, " already exists, do you wish to overwrite it?"), "y")

			if answ == "y" || answ == strings.ToLower("YES") {

				wg.Add(1)
				sem <- struct{}{}

				go func(srcPath, srcName, targetDir string) {
					defer func() {
						wg.Done()
						currentFileName <- fmt.Sprint("Written ", filepath.Join(targetDir, srcName), " ")
						<-sem
					}()

					os.Remove(filepath.Join(targetDir, srcName))
					writeFile(srcPath, srcName, targetDir)
				}(srcPath, srcName, targetDir)

			}

		}

	}
}

func Relocate(prompt string, deleteCopied bool, selected *navigator.Selected, nav *navigator.Navigator) {
	answ := WaitInput(fmt.Sprintf("%s %s", prompt, "selected into the current directory? :"), "y")

	if answ != "y" {
		return
	}

	go func() {
		terminal.PrintProgress(stopProgress, currentFileName)
	}()

	for entry := range selected.GetAll() {

		if *entry.Path == nav.CurrentPath {
			utils.ShowErrAndContinue(errors.New("copying / moving within same directory is not permitted"))
			return
		}

		if strings.HasPrefix(nav.CurrentPath, entry.FullPath()) {
			utils.ShowErrAndContinue(errors.New("cannot move / copy a folder into itself"))
			return
		}

		wg.Add(1)
		go func(path string, name string, currentPath string) {
			defer wg.Done()

			move(path, name, currentPath, deleteCopied)
		}(*entry.Path, entry.Name, nav.CurrentPath)

	}
	wg.Wait()

	stopProgress <- struct{}{}
	// wg.Wait()

	// TODO: after copying sort according to current sorting alghoritm
}
