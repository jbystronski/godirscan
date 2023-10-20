package controller

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
)

func (c *Controller) executeCmd() {
	input := c.wrapInput("Run command: ", "")

	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	//	cmd.Stdin = os.Stdin
	common.ClearScreen()
	fmt.Println("Press esc to return, command execution output: " + "\033[0m")
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		c.ErrorChan <- err
	}

	for {

		_, key, err := keyboard.GetKey()
		if err != nil {
			c.ErrorChan <- err
		}

		if key == keyboard.KeyEsc {
			break
		}

	}
	c.restoreScreen()
}
