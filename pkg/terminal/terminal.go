package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/navigator"
)

const (
	Space = " "

	Line = "\u2594"

	Bslash       = "\u2572"
	Fslash       = "\u2571"
	DbVert       = "\u2551"
	Texture      = "\u2591"
	ResetFmt     = "\033[0m"
	BoldFmt      = "\033[1m"
	DbHoriz      = "\u2550"
	CursorTop    = "\033[H"
	Hseparator   = "\u2500"
	CornerLine   = "\u2514"
	TeeLine      = "\u251c"
	Blink        = "\033[?25h"
	ClearLineFmt = "\033[2K"
	ReturnFmt    = "\r"
	CursorLeft   = "\033[D"
)

var Segment = strings.Join([]string{Fslash, Texture, Fslash}, "")

// const (
// 	EmptyIndent    = "   "
// 	VerticalSep    = "\u2502"
//
//
// 	HorizontalLine = "\u2500"
//
//

// 	FmtDir         = Bold + Blue
// 	FmtFile        = Yellow
// )

// boldFmt    = "\033[1m"
// italic     = "\033[3m"
// resetFmt   = "\033[0m"
// underscore = "\033[4m"
// black      = "\033[30m"

// tlCorner      = "\033(0\x6C\033(B"
// trCorner      = "\033(0\x71\033(B"
// blCorner      = "\033(0\x6D\033(B"
// brCorner      = "\033(0\x6A\033(B"
// lrDiag        = "\033(0\x2F\033(B"
// hTop          = "\033(0\x48\033(B"

// emptyIndent = "   "
// vSeparator  = "\u2502"
// teeLine     = "\u251c"
// cornerLine  = "\u2514"
// hSeparator  = "\u2500"
// dbHoriz     = "\u2550"
// line        = "\u2594"
// dbVert      = "\u2551"
// segment     = "\u2571" + "\u2591" + "\u2571"
// block       = "\u2586"
// space       = " "
// fSlash      = "\u2571"
// bSlash      = "\u2572"
//
// texture     = "\u2591"

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

func GetNumVisibleLines() int {
	rows, _ := getNumVisibleRowsAndCols()

	return rows
}

func getNumVisibleCols() int {
	_, cols := getNumVisibleRowsAndCols()

	return cols
}

func getNumVisibleRowsAndCols() (int, int) {
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
	split := strings.Split(st, " ")

	rows, _ := strconv.Atoi(split[0])
	cols, _ := strconv.Atoi(split[1])

	return rows, cols, nil
}

func MoveCursorTop() {
	fmt.Print(CursorTop)
}

func ClearLine() {
	fmt.Printf(ClearLineFmt)
}

func CarriageReturn() {
	fmt.Printf(ReturnFmt)
}

func RenderOutput(n *navigator.Navigator, s *navigator.Selected) {
	MoveCursorTop()
	n.NumVisibleLines = GetNumVisibleLines() - 2

	n.SetEndLine(n.GetStartLine() + n.NumVisibleLines)

	if n.GetEndline() > n.GetEntriesLength() {
		n.SetEndLine(n.GetEntriesLength())
	}

	var sep string
	ClearLine()

	PrintHeader(n.GetCurrentPath() + Space + entry.PrintSizeAsString(*n.GetDirSize()))

	if !n.HasEntries() {
		PrintEmpty()
	}
	for i := n.GetStartLine(); i < n.GetEndline(); i++ {

		if i == n.GetEntriesLength()-1 {
			sep = CornerLine + Hseparator
		} else {
			sep = TeeLine + Hseparator
		}
		ClearLine()
		if n.GetCurrentIndex() == i {
			HighlightRow(sep, *n.GetEntry(i))
		} else if _, ok := s.SelectedEntries[n.GetEntry(i)]; ok {
			MarkRow(sep, *n.GetEntry(i))
		} else {
			PrintRow(sep, *n.GetEntry(i))
		}

	}

	// stop := make(chan struct{})

	// ShowCursor() // Show cursor
}

func FlushOutput(n *navigator.Navigator, s *navigator.Selected) {
	ClearScreen()
	RenderOutput(n, s)
}

func ShowCursor() {
	fmt.Println(Blink)
}

func ResetFlushOutput(n *navigator.Navigator, s *navigator.Selected) {
	ClearScreen()
	n.Reset()
	RenderOutput(n, s)
}
