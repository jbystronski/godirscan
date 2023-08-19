package task

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/terminal"
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
	fmt.Println(terminal.CurrentTheme.Accent + "Ppress esc to return, command execution output: " + terminal.ResetFmt)
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return
	}

	// _ = k.Open()

	// defer func() {
	// 	k.Close()
	// }()

	for {

		_, key, err := k.GetKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if key == k.KeyEsc {
			return
		}

	}
}

func ExecCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	//	cmd.Stdin = os.Stdin
	terminal.ClearScreen()
	fmt.Println(terminal.CurrentTheme.Accent + "Press esc to return, command execution output: " + terminal.ResetFmt)
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return
	}
	// _ = k.Open()

	// defer func() {
	// 	k.Close()
	// }()

	for {

		_, key, err := k.GetKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if key == k.KeyEsc {
			//	terminal.ClearScreen()
			return
		}

	}
}
