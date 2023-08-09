package task

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/terminal"
	"github.com/jbystronski/godirscan/pkg/utils"
)

func ExecuteDefault(path string) {
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
	terminal.ClearScreen()
	fmt.Println(terminal.CurrentTheme.Accent + "Command execution output (press esc to return): " + terminal.ResetFmt)
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		utils.ShowErrAndContinue(err)
		return
	}

	for {

		_, key, err := k.GetKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if key == k.KeyEsc {
			terminal.ClearScreen()
			return
		}

	}
}

func ExecCommand(input string) {
	args := strings.Fields(input)

	cmd := exec.Command(args[0], args[1:]...)
	//	cmd.Stdin = os.Stdin
	terminal.ClearScreen()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		utils.ShowErrAndContinue(err)
		return
	}
	fmt.Println(terminal.CurrentTheme.Accent + "Command output (press esc to return): " + terminal.ResetFmt)
	fmt.Println()
	for {

		_, key, err := k.GetKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if key == k.KeyEsc {
			terminal.ClearScreen()
			return
		}

	}
}
