package task

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/navigator"
	"github.com/jbystronski/godirscan/pkg/terminal"
	"github.com/jbystronski/godirscan/pkg/utils"
)

func DeleteSelected(selected *navigator.Selected, nav *navigator.Navigator) (ok bool, err error) {
	answ, inputErr := WaitInput("Delete selected entries", "y")
	if inputErr != nil {
		err = inputErr
		return
	}

	if answ == "y" {

		stopProgress := make(chan struct{}, 1)
		currentFileName := make(chan string)

		go func() {
			terminal.PrintProgress(stopProgress, currentFileName)
		}()

		var deleteGroup sync.WaitGroup

		for key := range selected.GetAll() {

			deleteGroup.Add(1)
			go func(en *entry.Entry) {
				defer func() {
					deleteGroup.Done()
					currentFileName <- fmt.Sprint("Deleted ", en.Name, " ")
				}()

				removeErr := os.RemoveAll(en.FullPath())
				if removeErr != nil {
					if errors.Is(removeErr, os.ErrPermission) {
						utils.ShowErrAndContinue(err)
					} else {
						err = removeErr
						return

					}
				}
			}(key)
			// time.Sleep(time.Millisecond * 50)
		}
		deleteGroup.Wait()

		stopProgress <- struct{}{}
		ok = true

	}
	return
}
