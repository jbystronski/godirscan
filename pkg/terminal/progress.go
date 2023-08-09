package terminal

import (
	"fmt"
	"strings"
	"time"
)

func PrintProgress(stop <-chan struct{}, msg <-chan string) {
	//	stages := [...]string{"|", Fslash, Hseparator, Bslash}
	stages := [...]string{".   ", "..  ", "... ", "...."}

	// testNames := []string{"some file", "translations.pdf", "windowsKey29.zip", "codeV1.2xsafdas.deb"}

	var currentMsg string

	var currStage uint8
	output := ""
	showProgress := true

	for showProgress {
		select {
		case <-stop:
			showProgress = false
			return
		// case <-time.After(time.Second * 5):
		// 	showProgress = false
		// 	return

		case receivedMsg := <-msg:

			// fmt.Println(receivedMsg)
			// time.Sleep(time.Millisecond * 100)
			currentMsg = receivedMsg

		default:
			output += currentMsg + stages[currStage]

			if currStage == uint8(len(stages)-1) {
				currStage = 0
			} else {
				currStage++
			}

			// fmt.Print(len(output) - 1)
			// time.Sleep(time.Second * 2)
			ClearLine()
			CarriageReturn()
			fmt.Print(output)
			time.Sleep(time.Millisecond * 100)
			fmt.Print(strings.Repeat(CursorLeft, len(output)))
			output = ""

		}
	}
}
