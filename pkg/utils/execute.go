package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func executeDefault(path string, errHandle func(error)) {
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		errHandle(err)
	}
}

func execCommand(input string, errHandle func(error)) {
	args := strings.Fields(input)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		errHandle(err)
	}
}
