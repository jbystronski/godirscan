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

func DeleteSelected(selected *navigator.Selected, nav *navigator.Navigator) {
	answ := WaitInput("Delete selected entries", "y")

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

				err := os.RemoveAll(en.FullPath())
				if err != nil {
					if errors.Is(err, os.ErrPermission) {
						utils.ShowErrAndContinue(err)
						return
					} else {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}(key)
			// time.Sleep(time.Millisecond * 50)
		}
		deleteGroup.Wait()

		stopProgress <- struct{}{}

	}
}
