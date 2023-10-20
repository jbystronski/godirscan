package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	space = " "

	Line = "\u2594"

	backslash           = "\u2572"
	slash               = "\u2571"
	doubleVertical      = "\u2551"
	texture             = "\u2591"
	reset               = "\033[0m"
	bold                = "\033[1m"
	italic              = "\033[3m"
	underscore          = "\033[4m"
	doubleHorizontal    = "\u2550"
	cursorTop           = "\033[H"
	horizontalSeparator = "\u2500"
	corner              = "\u2514"
	tee                 = "\u251c"
	ShowCursor          = "\033[?25h"
	hideCursor          = "\033[?25l"
	clearLineFormat     = "\033[2K"
	returnFormat        = "\r"
	cursorLeft          = "\033[D"
	cursorRight         = "\033[C"

	borderHorizontal    = "\u2501"
	borderHorizontalAlt = "\u2550"
	borderVertical      = "\u2503"
	borderVerticalAlt   = "\u2551"
	// topLeftBorderCorner        = "\u250F"
	// topLeftBorderCornerAlt     = "\u2554"
	// topRightBorderCorner       = "\u2513"
	// topRightBorderCornerAlt    = "\u2557"
	// bottomLeftBorderCorner     = "\u2517"
	// bottomLeftBorderCornerAlt  = "\u255A"
	// bottomRightBorderCorner    = "\u251B"
	// bottomRightBorderCornerAlt = "\u255D"
)

var segment = strings.Join([]string{slash, texture, slash}, "")

func r(s string, times int) string {
	return strings.Repeat(s, times)
}

func HideCursor() {
	fmt.Print(hideCursor)
}

func PrintBanner(t Theme) {
	fmt.Print(t.Main)

	fmt.Print("\u2554" + r(doubleHorizontal, 90) + "\u2557\n")
	fmt.Print(doubleVertical + r(space, 90) + doubleVertical + "\n")

	fmt.Print(doubleVertical + r(" ", 6) + slash + r(texture, 8) + slash + r(space, 74) + doubleVertical + "\n")

	fmt.Print(doubleVertical + r(" ", 5) + segment + r(Line, 6) + r(space, 2) + slash + r(texture, 7) + slash + r(space, 1) + slash + r(texture, 6) + backslash + r(space, 2) + slash + r(texture, 1) + slash + r(space, 1) + slash + r(texture, 6) + backslash + r(space, 2) + slash + r(texture, 6) + slash + r(space, 1) + slash + r(texture, 6) + slash + r(space, 1) + slash + r(texture, 7) + slash + r(space, 1) + slash + r(texture, 2) + backslash + r(space, 2) + slash + r(texture, 1) + slash + r(space, 3) + doubleVertical + "\n")
	fmt.Print(doubleVertical + r(" ", 4) + segment + r(space, 2) + slash + r(texture, 3) + slash + r(space, 1) + segment + r(Line, 3) + segment + space + segment + r(Line, 3) + segment + r(space, 1) + segment + space + segment + r(Line, 3) + segment + r(space, 1) + slash + r(texture, 5) + backslash + r(space, 2) + segment + r(Line, 5) + r(space, 1) + segment + r(Line, 3) + segment + r(space, 1) + segment + backslash + texture + backslash + segment + r(space, 4) + doubleVertical + "\n")
	fmt.Print(doubleVertical + r(" ", 3) + segment + r(space, 3) + r(Line, 1) + segment + space + segment + r(space, 3) + segment + space + segment + r(space, 3) + segment + r(space, 1) + segment + space + slash + r(texture, 7) + r(Line, 0) + r(space, 0) + slash + r(space, 2) + r(Line, 4) + slash + texture + slash + space + segment + r(space, 6) + slash + r(texture, 7) + slash + space + segment + r(space, 2) + backslash + r(texture, 2) + slash + r(space, 5) + doubleVertical + "\n")
	fmt.Print(doubleVertical + r(" ", 2) + slash + r(texture, 8) + slash + space + slash + r(texture, 7) + slash + space + slash + r(texture, 7) + slash + space + slash + texture + slash + r(space, 1) + segment + r(Line, 4) + backslash + texture + backslash + r(space, 1) + slash + r(texture, 5) + slash + space + slash + r(texture, 6) + slash + space + segment + r(Line, 3) + segment + space + segment + r(space, 3) + segment + r(space, 6) + doubleVertical + "\n")

	fmt.Print(doubleVertical + r(" ", 2) + r(Line, 9) + r(space, 2) + r(Line, 8) + r(space, 2) + r(Line, 8) + r(space, 2) + Line + Line + r(space, 2) + r(Line, 2) + r(space, 6) + r(Line, 2) + space + r(Line, 6) + r(space, 2) + r(Line, 7) + r(space, 2) + r(Line, 2) + r(space, 4) + r(Line, 2) + r(space, 2) + r(Line, 2) + r(space, 4) + r(Line, 2) + r(space, 7) + doubleVertical + "\n")

	fmt.Print("\u255A" + r(doubleHorizontal, 90) + "\u255D")

	fmt.Print(reset)

	fmt.Print("\n\n")
}

func PrintHelpers(t Theme) {
	helper := func(s string) {
		fmt.Printf("%v%v%s%v%s%v%v%s", bold, t.BgHeader, space, ColorsMap["bright_white"], s, space, reset, space)
	}

	const (
		op1 string = "F2 - Help"
		op2        = "Ctrl M - Match"
		op3        = "Esc - Quit"
	)

	helper(op1)
	helper(op2)
	helper(op3)
}

// func PrintHelp(t theme.Theme, p filesystem.Controller) {
// 	//for i, s := range selected {
// 	//	fmt.Printf("%d - %s", i, s)
// 	//}

// 	PrintBanner(t)
// 	fmt.Println("total lines: ", p.GetTotalLines())
// 	fmt.Println("terminal columns: ", GetNumVisibleCols())

// 	fmt.Println("output lines: ", p.GetLines())
// 	fmt.Println("output first line: ", p.GetOutputFirstLine())
// 	fmt.Println("output last line: ", p.GetOutputLastLine())
// 	fmt.Println("current index ", p.GetIndex())

// 	fmt.Println("first index: ", p.GetChunkStart())
// 	fmt.Println("last index:", p.GetChunkEnd())

// 	fmt.Println("Path of selected ", p.Selected.SelectedEntriesPath)

// 	for k := range p.Selected.GetAll() {
// 		fmt.Println(k)
// 	}

// 	for {
// 		_, key, err := keyboard.GetKey()
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		if key == keyboard.KeyEsc {
// 			return
// 		}

// 	}

// 	// printOptions(optionsList)
// }

func PrintPrompt(prompt string, t Theme) string {
	return fmt.Sprintf("%s%s %s %s", t.BgPrompt, t.Prompt+bold, prompt, reset)
}

func Cell(rowNumber, colNumber int) {
	fmt.Printf("\033[%d;%dH", rowNumber, colNumber)
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

func ClearRow(row, offsetLeft, length int) {
	Cell(row, offsetLeft)
	fmt.Print(strings.Repeat(space, length))
	Cell(row, offsetLeft)
}

func ClearLine() {
	fmt.Printf(clearLineFormat)
}

func CarriageReturn() {
	fmt.Printf(returnFormat)
}

func Erase(offsetTop, offsetLeft, cells int) {
	Cell(offsetTop, offsetLeft)
	for cells != 0 {
		fmt.Print(space)
		cells--
	}
	Cell(offsetTop, offsetLeft)
}

func PrintProgress(doneChan <-chan struct{}, messageChan <-chan string, offsetLeft, offsetTop int) {
	defer func() {
		Cell(offsetTop, offsetLeft)

		ClearLine()
		CarriageReturn()
	}()

	Cell(offsetTop, offsetLeft)

	stages := [...]string{".   ", "..  ", "... ", "...."}

	var currStage uint8
	output := ""
	showProgress := true

	for showProgress {
		select {
		case <-doneChan:
			showProgress = false

			return

		default:
			msg := <-messageChan

			output += msg + stages[currStage]

			if currStage == uint8(len(stages)-1) {
				currStage = 0
			} else {
				currStage++
			}
			Cell(offsetTop, offsetLeft)

			ClearLine()
			CarriageReturn()

			fmt.Print(output)

			//	time.Sleep(time.Millisecond * 100)
			// fmt.Print(strings.Repeat(CursorLeft, len(output)))
			//	time.Sleep(time.Millisecond * 40)
			output = ""
			ClearLine()
			CarriageReturn()

		}
	}
}

func NumVisibleLines() int {
	rows, _ := NumVisibleRowsAndCols()

	return rows
}

func NumVisibleCols() int {
	_, cols := NumVisibleRowsAndCols()

	return cols
}

func NumVisibleRowsAndCols() (int, int) {
	rows, cols, _ := terminalSize()
	return rows, cols
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

	rows, cols := 0, 0

	switch runtime.GOOS {
	case "windows":
		// Parse rows and columns from the output of 'mode' command on Windows
		rows, _ = strconv.Atoi(strings.Split(st, "\n")[2][8:])
		cols, _ = strconv.Atoi(strings.Split(st, "\n")[3][9:])
	default:
		// Parse rows and columns from the output of 'stty size' command on Unix-like systems
		split := strings.Fields(st)
		if len(split) >= 2 {
			rows, _ = strconv.Atoi(split[0])
			cols, _ = strconv.Atoi(split[1])
		}
	}

	return rows, cols, nil
}

func PaneWidth() int {
	borderWidth := 1

	baseWidth := NumVisibleCols()
	if baseWidth%2 != 0 {
		baseWidth -= 1
	}

	return (baseWidth - 2*borderWidth) / 2
}

func FlashError(err error) {
	numOfLines := NumVisibleLines()
	numOfCols := NumVisibleCols()

	Erase(numOfLines, 1, numOfCols)
	fmt.Print(err)
	time.Sleep(time.Second * 2)
	Erase(numOfLines, 1, numOfCols)
	Cell(numOfLines, 2)
	PrintHelpers(*CurrentTheme)
}

func GetRootDirectory() string {
	wd, _ := os.Getwd()

	return filepath.VolumeName(wd) + string(filepath.Separator)
}
