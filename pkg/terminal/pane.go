package terminal

import (
	"fmt"
	"strings"
)

const (
	borderWidth   = 1
	reserverdRows = 4
	startRow      = 3
)

func PrintPane(line, startCol, endCol int) {
	//	height := GetNumVisibleLines() - reserverdRows

	// ClearScreen()
	Cell(line, startCol)
	printTopBorder()
	// Cell(8, 1)
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
	// fmt.Print(colNumber)
}

func GetPaneWidth() int {
	return (getNumVisibleCols() - 2*borderWidth) / 2
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
