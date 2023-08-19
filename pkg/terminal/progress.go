package terminal

import (
	"fmt"
)

func PrintProgress(stop <-chan struct{}, msg <-chan string) {
	//	stages := [...]string{"|", Fslash, Hseparator, Bslash}
	stages := [...]string{".   ", "..  ", "... ", "...."}

	var currentMsg string

	var currStage uint8
	output := ""
	showProgress := true

	for showProgress {
		select {
		case <-stop:
			showProgress = false
			return

		case receivedMsg := <-msg:

			currentMsg = receivedMsg

		default:
			output += currentMsg + stages[currStage]

			if currStage == uint8(len(stages)-1) {
				currStage = 0
			} else {
				currStage++
			}

			Cell(totalLines, 1)
			// ClearLine()
			// CarriageReturn()
			fmt.Print(output)
			//	time.Sleep(time.Millisecond * 100)
			// fmt.Print(strings.Repeat(CursorLeft, len(output)))
			output = ""

		}
	}
}
