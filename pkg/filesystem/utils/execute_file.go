package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
)

func ExecuteFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		{
			cmd = exec.Command("open", path)
		}
	case "windows":
		{
			cmd = exec.Command("cmd", "/c", "start", path)
		}
	default:
		{
			cmd = exec.Command("xdg-open", path)
		}
	}
	common.ClearScreen()
	fmt.Println(common.CurrentTheme.Accent + "Press esc to return, command execution output: " + "\033[0m")
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	for {

		_, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		if key == keyboard.KeyEsc {
			return nil
		}

	}
}
