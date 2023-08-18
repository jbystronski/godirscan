package terminal

import (
	"fmt"
	"strings"
)

const (
	borderWidth   = 1
	promptLines   = 1
	headerLines   = 1
	reserverdRows = 4
)

var paneWidth, OutputFirstLine, OutputLastLine, OutputLines, PromptLine int

func PrintPane(line, startCol, endCol int) {
	Cell(line, startCol)
	printTopBorder()

	for i := 3; i < GetNumVisibleLines(); i++ {
		Cell(i, startCol)
		printBorder()

		Cell(i, endCol)
		printBorder()
		fmt.Print("\n")
		line++
	}
	Cell(line, startCol)
	printBottomBorder()
}

func GetPaneWidth() int {
	return paneWidth
}

func SetLayout() {
	paneWidth = (GetNumVisibleCols() - 2*borderWidth) / 2
	OutputFirstLine = headerLines + borderWidth + 1
	OutputLastLine = GetNumVisibleLines() - promptLines - borderWidth
	OutputLines = OutputLastLine - OutputFirstLine
	PromptLine = OutputLastLine + borderWidth + 1
}

func printBorder() {
	fmt.Printf("%v%s%v", CurrentTheme.Main, borderVertical, ResetFmt)
}

func printTopBorder() {
	fmt.Printf("%v%s%s%s%v\n", CurrentTheme.Main, topLeftBorderCorner, strings.Repeat(borderHorizontal, GetPaneWidth()-2), topRightBorderCorner, ResetFmt)
}

func printBottomBorder() {
	fmt.Printf("%v%s%s%s%v\n", CurrentTheme.Main, bottomLeftBorderCorner, strings.Repeat(borderHorizontal, GetPaneWidth()-2), bottomRightBorderCorner, ResetFmt)
}
