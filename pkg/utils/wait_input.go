package utils

import (
	"fmt"
	"time"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func waitUserInput(prompt, placeholder string, forwardOutput func(string)) {
	if !k.IsStarted(time.Millisecond * 100) {
		k.Open()

		defer k.Close()
	}

	print := func(s string) {
		terminal.CarriageReturn()
		fmt.Print(fPrompt(prompt) + " " + s)
	}

	output := placeholder

	terminal.ClearLine()
	print(output)
	for {

		char, key, err := k.GetKey()

		printDefaultErrorAndExit(err)

		if key == quitKey {

			terminal.ClearLine()
			terminal.CarriageReturn()
			break

		} else if key == k.KeyEnter {
			terminal.ClearLine()
			terminal.CarriageReturn()
			forwardOutput(output)
			break

		} else if key == k.KeyBackspace || key == k.KeyBackspace2 {
			if len(output) > 0 {
				output = output[:len(output)-1]
				terminal.ClearLine()

				print(output)
			}
		} else if key == k.KeySpace {
			output += " "
			print(output)
		} else if char != 0 {
			output += string(char)

			print(output)
		}

	}
}
