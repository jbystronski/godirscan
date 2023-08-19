package terminal

import (
	"fmt"
	"strings"
)

const (
	borderWidth      = 1
	helperLineHeight = 1
	promptLineHeight = 1
	headerLineHeight = 1
	reserverdRows    = 4
)

var paneWidth, OutputWidth, helperLine, topBorderLine, bottomBorderLine, OutputFirstLine, OutputLastLine, OutputLines, PromptLine, totalLines int

func printLeftPane() {
	PrintPane(1, paneWidth)
}

func printRightPane() {
	PrintPane(paneWidth+1, paneWidth*2)
}

func PrintPane(startCol, endCol int) {
	Cell(topBorderLine, startCol)
	printTopBorder()
	Cell(bottomBorderLine, startCol)
	printBottomBorder()

	currentRow := OutputFirstLine

	for currentRow <= OutputLastLine {
		Cell(currentRow, startCol)
		printBorder()

		Cell(currentRow, endCol)
		printBorder()

		currentRow++
	}
	Cell(currentRow, startCol)
	printBorder()
	Cell(currentRow, endCol)
	printBorder()
	// Cell(currentRow, startCol+1)
	// fmt.Print("Prompt")
}

func PrintPanes() {
	printLeftPane()
	printRightPane()
	Cell(helperLine, 1)
	printHelpers()
}

func GetPaneWidth() int {
	return paneWidth
}

func SetLayout() {
	totalLines = GetNumVisibleLines()
	paneWidth = (GetNumVisibleCols() - 2*borderWidth) / 2

	OutputFirstLine = headerLineHeight + borderWidth + 1
	OutputLastLine = totalLines - helperLineHeight - borderWidth - promptLineHeight

	//	OutputLines = TotalLines - headerLineHeight - borderWidth - promptLineHeight - borderWidth - helperLineHeight
	OutputLines = OutputLastLine - OutputFirstLine

	PromptLine = totalLines - 2
	topBorderLine = 2
	bottomBorderLine = totalLines - 1
	OutputWidth = paneWidth - 2
	helperLine = totalLines
}

func printBorder() {
	fmt.Printf("%v%s%v", CurrentTheme.Main, borderVertical, ResetFmt)
}

func printTopBorder() {
	fmt.Printf("%v%s%s%s%v\n", CurrentTheme.Main, topLeftBorderCorner, strings.Repeat(borderHorizontal, paneWidth-2), topRightBorderCorner, ResetFmt)
}

func printBottomBorder() {
	fmt.Printf("%v%s%s%s%v\n", CurrentTheme.Main, bottomLeftBorderCorner, strings.Repeat(borderHorizontal, paneWidth-2), bottomRightBorderCorner, ResetFmt)
}
