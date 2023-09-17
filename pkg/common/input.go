package common

import (
	"fmt"

	k "github.com/eiannone/keyboard"
)

func WaitInput(prompt, output string, coords Coords, errChan chan<- error,
) (result string) {
	defer func() {
		Cell(NumVisibleLines(), 1)
		fmt.Print("\033[?25l")
	}()

	Cell(coords.Y, coords.X)
	ClearRow(coords.Y, coords.X, PaneWidth()-2)

	fmt.Print(PrintPrompt(prompt, *CurrentTheme) + " ")

	cmdLine := NewCommandLine(coords.Y, coords.X+len(prompt)+1, PaneWidth()-2, output)

	for {

		char, key, getKeyErr := k.GetKey()
		if getKeyErr != nil {
			errChan <- getKeyErr
			return

		}

		switch key {
		case k.KeyEsc:
			ClearLine()
			CarriageReturn()
			return

		case k.KeyEnter:
			ClearLine()
			CarriageReturn()
			result = cmdLine.GetOutput()

			return

		case k.KeyArrowRight:

			cmdLine.NextCol()

		case k.KeyArrowLeft:
			cmdLine.PrevCol()

		case k.KeyBackspace, k.KeyBackspace2:
			cmdLine.Backspace()

		case k.KeyDelete:

			cmdLine.DeleteChar()

		case k.KeyTab:
			for i := 0; i <= 3; i++ {
				cmdLine.InsertChar(' ')
			}

		case k.KeySpace:

			cmdLine.InsertChar(' ')

		case k.KeyHome:

			cmdLine.GoToLineStart()

		case k.KeyEnd:
			cmdLine.GoToLineEnd()

		default:

			if char != 0 {
				cmdLine.InsertChar(char)
			}

		}

	}
}
