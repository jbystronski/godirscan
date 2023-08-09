package utils

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func PrintDefaultErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func ShowErrAndContinue(err error) {
	if err != nil {
		// blockNavigation = true
		terminal.ClearScreen()
		fmt.Printf("%s%s\n\n%s%s", terminal.CurrentTheme.Accent, err, terminal.ResetFmt, "press space to continue")

		for {

			_, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}

			if key == keyboard.KeySpace {
				// blockNavigation = false

				break
			}

		}

	}
}
