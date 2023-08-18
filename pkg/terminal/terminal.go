package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/navigator"
	"golang.org/x/text/unicode/norm"
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
	cursorRight  = "\033[C"
	// borderHorizontal        = "\u2550"
	// borderVertical          = "\u2551"
	// topLeftBorderCorner     = "\u2554"
	// topRightBorderCorner    = "\u2557"
	// bottomLeftBorderCorner  = "\u255A"
	// bottomRightBorderCorner = "\u255D"
	borderHorizontal        = "\u2501"
	borderVertical          = "\u2503"
	topLeftBorderCorner     = "\u250F"
	topRightBorderCorner    = "\u2513"
	bottomLeftBorderCorner  = "\u2517"
	bottomRightBorderCorner = "\u251B"
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

func moveToEndOfLine() {
	fmt.Print("\033[K")
}

func Cell(rowNumber, colNumber int) {
	fmt.Printf("\033[%d;%dH", rowNumber, colNumber)
}

func trimLine(line string, max int) string {
	if utf8.RuneCountInString(line) > max {
		line = line[0:max]
		//*line += *line + "..."
	}
	return fmt.Sprintf("%v%v", line, ResetFmt)
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

func GetNumVisibleLines() int {
	rows, _ := getNumVisibleRowsAndCols()

	return rows
}

func GetNumVisibleCols() int {
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

func MoveCursorTop() {
	fmt.Print(CursorTop)
}

func ClearRow(row, offsetLeft, length int) {
	Cell(row, offsetLeft)
	fmt.Print(strings.Repeat(Space, length))
}

func ClearLine() {
	fmt.Printf(ClearLineFmt)
}

func CarriageReturn() {
	fmt.Printf(ReturnFmt)
}

func FlushOutput(n *navigator.Navigator, s *navigator.Selected) {
	//	ClearScreen()
	RenderOutput(n, s)
}

func ResetFlushOutput(n *navigator.Navigator, s *navigator.Selected) {
	// ClearScreen()
	Cell(1, 1)
	printHeader(n.CurrentPath + Space + entry.PrintSizeAsString(*n.GetDirSize()))

	//	ClearScreen()
	n.Reset()
	RenderOutput(n, s)
}

func RenderOutput(n *navigator.Navigator, s *navigator.Selected) {
	ClearRow(1, n.StartCell, n.RowWidth)
	Cell(1, n.StartCell-1)

	st := printHeader(n.CurrentPath + Space + entry.PrintSizeAsString(*n.GetDirSize()))

	fmt.Print(st)
	// totalLines := GetNumVisibleLines()

	currentRow := OutputFirstLine
	// lastRow := totalLines - 2
	// 23
	// visibleRows := lastRow - currentRow

	var startIndex, endIndex int

	// 34 > 23
	if n.CurrentIndex > OutputLines {
		startIndex = n.CurrentIndex - OutputLines
		endIndex = n.CurrentIndex
	} else {
		startIndex = 0
		endIndex = OutputLastLine - OutputFirstLine
	}

	var sep string
	for i := startIndex; i <= endIndex; i++ {

		ClearRow(currentRow, n.StartCell, n.RowWidth)
		Cell(currentRow, n.StartCell)
		if i >= n.GetEntriesLength() {
			currentRow++
			continue
		} else {

			if i == n.GetEntriesLength()-1 {
				sep = CornerLine + Hseparator
			} else {
				sep = TeeLine + Hseparator
			}
			var output string
			if n.CurrentIndex == i && n.IsActive {
				output = highlightRow(sep, *n.GetEntry(i))
			} else if _, ok := s.SelectedEntries[n.GetEntry(i)]; ok {
				output = MarkRow(sep, *n.GetEntry(i))
			} else {
				output = printRow(sep, *n.GetEntry(i))
			}

			fmt.Print(trimLine(output, n.RowWidth))

			currentRow++
		}

	}
	Cell(PromptLine, 1)
	// Cell(GetNumVisibleLines(), 1)
	// MoveCursorTop()
	// n.NumVisibleLines = GetNumVisibleLines() - reserverdRows

	// n.SetEndLine(n.GetStartLine() + n.NumVisibleLines)

	// PrintHeader(n.CurrentPath + Space + entry.PrintSizeAsString(*n.GetDirSize()))
	// printTopBorder()

	// if n.EndLine > n.GetEntriesLength() {
	// 	n.SetEndLine(n.GetEntriesLength())
	// }

	// var sep string
	// ClearLine()

	// if !n.HasEntries() {
	// 	PrintEmpty()
	// }
	// for i := n.StartLine; i < n.EndLine; i++ {

	// 	if i == n.GetEntriesLength()-1 {
	// 		sep = CornerLine + Hseparator
	// 	} else {
	// 		sep = TeeLine + Hseparator
	// 	}
	// 	ClearLine()
	// 	if n.CurrentIndex == i {
	// 		HighlightRow(sep, *n.GetEntry(i))
	// 	} else if _, ok := s.SelectedEntries[n.GetEntry(i)]; ok {
	// 		MarkRow(sep, *n.GetEntry(i))
	// 	} else {
	// 		PrintRow(sep, *n.GetEntry(i))
	// 	}

	// }
}

func ShowCursor() {
	fmt.Println(Blink)
}

func MoveCursorLeft(times int) {
	fmt.Print(strings.Repeat(CursorLeft, times))
}

func MoveCursorRight(times int) {
	fmt.Print(strings.Repeat(cursorRight, times))
}

func RuneToUtf8String(r rune) string {
	if r < utf8.RuneSelf && r >= 0 {
		return string(r)
	}

	norms := []norm.Form{norm.NFC, norm.NFKC, norm.NFD, norm.NFKD}

	var normalized string

	for _, form := range norms {
		normalized = form.String(string(r))
		fmt.Printf(" %v (%v bytes)\n", normalized, len(normalized))
		time.Sleep(time.Second * 1)
		if len(normalized) == 1 {
			return normalized
		}
	}

	return string(r)

	// normalized := norm.NFKC.String(string(r))

	// displayWidth := utf8.RuneCountInString(normalized)

	// if displayWidth > 1 {

	// 	truncated := normalized[:utf8.RuneLen([]rune(normalized)[0])]
	// 	return truncated
	// }

	// return normalized
}
