package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

func GetNumVisibleLines() int {
	rows, _ := GetNumVisibleRowsAndCols()

	return rows
}

func GetNumVisibleCols() int {
	_, cols := GetNumVisibleRowsAndCols()

	return cols
}

func GetNumVisibleRowsAndCols() (int, int) {
	rows, cols, _ := terminalSize()
	return rows, cols
}

func GetUserDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func ResolveUserDirectory(fPath *string) {
	switch runtime.GOOS {
	case "darwin":
		{
			break
		}
	case "windows":
		{
			break
		}
	default:
		{
			if strings.HasPrefix(*fPath, "~") {
				currentUser, err := user.Current()
				if err != nil {
					panic(err)
				}

				*fPath = strings.Replace(*fPath, "~", currentUser.HomeDir, 1)

			}
		}
	}
}

func terminalSize() (int, int, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "windows":
		{
			cmd = exec.Command("cmd", "/c", "mode")
		}
	default:
		{
			cmd = exec.Command("stty", "size")
		}
	}

	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	st := string(out)
	split := strings.Split(st, " ")

	rows, _ := strconv.Atoi(split[0])
	cols, _ := strconv.Atoi(split[1])

	return rows, cols, nil
}

func MoveCursorTop() {
	fmt.Print(CursorTop)
}

func ClearScreen() {
	clearCommand := ""

	switch runtime.GOOS {
	case "linux", "darwin":
		clearCommand = "clear"
	case "windows":
		clearCommand = "cls"
	default:
		panic("Unsupported platform")
	}

	cmd := exec.Command(clearCommand)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ClearLine() {
	fmt.Printf("\033[2K")
}

func CarriageReturn() {
	fmt.Printf("\r")
}
