package task

import (
	"fmt"
	"os"
	"strings"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func WaitUserInput(prompt, placeholder string, forwardOutput func(string)) {
	print := func(s string) {
		terminal.CarriageReturn()
		fmt.Print(terminal.Prompt(prompt) + " " + s)
	}

	output := placeholder

	terminal.ClearLine()
	print(output)
	for {

		char, key, err := k.GetKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if key == k.KeyEsc {

			terminal.ClearLine()
			terminal.CarriageReturn()
			return

		} else if key == k.KeyEnter {
			terminal.ClearLine()
			terminal.CarriageReturn()
			output = strings.TrimSpace(output)
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
