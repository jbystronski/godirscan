package filesystem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/jbystronski/godirscan/pkg/common"
)

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
		//	currentFileName <- fmt.Sprint("Created symlink ", targetPath, " ")
		return true, nil

	}
	return false, nil
}

func writeFile(srcPath, targetPath string) error {
	srcFile, err := os.Open(srcPath)

	defer srcFile.Close()

	if err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)

	defer targetFile.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(srcFile, targetFile)

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

func move(srcPath, targetPath string, msgChan chan<- string, errChan chan<- error, sem chan struct{}, wg *sync.WaitGroup) {
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
				errChan <- err
				return
			}
			<-sem

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
				errChan <- err
				return
			}
			dc, readErr := os.ReadDir(srcPath)
			if readErr != nil {
				errChan <- readErr
				return
			}

			for _, entry := range dc {
				move(filepath.Join(srcPath, entry.Name()), filepath.Join(targetPath, entry.Name()), msgChan, errChan, sem, wg)
				if err != nil {
					errChan <- err
					return
				}

			}

		} else {
			wg.Add(1)

			go func(srcPath, targetPath string) {
				defer func(path string) {
					wg.Done()
					msgChan <- fmt.Sprint("Finished copying ", path)
					//	currentFileName <- fmt.Sprint("Finished writing ", path, " ")
				}(targetPath)
				//	currentFileName <- fmt.Sprint("Writing ", targetPath, " ")
				err = writeFile(srcPath, targetPath)

				if err != nil {
					errChan <- err
					return
				}

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

func relocate(action common.ControllerAction, errChan chan<- error, paths map[string]struct{}, targetDir string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, common.Cfg.MaxWorkers)

	doneChan := make(chan struct{})
	messageChan := make(chan string)

	go func() {
		common.PrintProgress(doneChan, messageChan, 1, common.NumVisibleLines())
	}()

	for srcPath := range paths {
		dir, file := filepath.Split(srcPath)

		if dir == targetDir {
			errChan <- errors.New("copying / moving within same directory is not permitted")
		}

		if strings.HasPrefix(targetDir, dir) {
			errChan <- errors.New("cannot move / copy a folder into itself")
		}

		targetPath := addSuffixIfExists(filepath.Join(targetDir, file))

		move(srcPath, targetPath, messageChan, errChan, sem, &wg)

		if action == common.Move {
			errChan <- os.RemoveAll(srcPath)
		}

	}

	wg.Wait()
}
