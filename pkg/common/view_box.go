package common

import (
	"fmt"
	"strings"
)

const (
	horizontalBorder        = "\u2501"
	verticalBorder          = "\u2503"
	topLeftBorderCorner     = "\u250F"
	topRightBorderCorner    = "\u2513"
	bottomLeftBorderCorner  = "\u2517"
	bottomRightBorderCorner = "\u251B"
	// horizontalSeparator     = "\u2500"
	// corner                  = "\u2514"
	// tee                     = "\u251c"
	// reset                   = "\033[0m"
	// bold                    = "\033[1m"
)

type ViewBox struct {
	offsetTopStart,
	offsetLeftStart,
	offsetTop,
	offsetBottom,
	totalLines,
	width,
	padding int
	theme           *Theme
	backgroundColor string
}

func (v *ViewBox) Theme() *Theme {
	return v.theme
}

func (v *ViewBox) SetTheme(t *Theme) {
	v.theme = t
}

func (v *ViewBox) Padding() int {
	return v.padding
}

func (v *ViewBox) SetPadding(p int) {
	v.padding = p
}

func (v *ViewBox) OffsetTopStart() int {
	return v.offsetTopStart
}

func (v *ViewBox) SetOffsetTopStart(o int) {
	v.offsetTopStart = o
}

func (v *ViewBox) OffsetLeftStart() int {
	return v.offsetLeftStart
}

func (v *ViewBox) SetOffsetLeftStart(o int) {
	v.offsetLeftStart = o
}

func (v *ViewBox) OffsetTop() int {
	return v.offsetTop
}

func (v *ViewBox) SetOffsetTop(o int) {
	v.offsetTop = o
}

func (v *ViewBox) OffsetBottom() int {
	return v.offsetBottom
}

func (v *ViewBox) SetOffsetBottom(o int) {
	v.offsetBottom = o
}

func (v *ViewBox) OutputFirstLine() int {
	return v.OffsetTop() + 1 + v.Padding()
}

func (v *ViewBox) TotalLines() int {
	return v.totalLines
}

func (v *ViewBox) SetTotalLines(tl int) {
	v.totalLines = tl
}

func (v *ViewBox) OutputLastLine() int {
	return v.TotalLines() - v.OffsetBottom()
}

func (v *ViewBox) Lines() int {
	return v.TotalLines() - v.OffsetBottom() - v.OffsetTop()
}

func (v *ViewBox) ContentLineStart() int {
	return v.OffsetLeftStart() + v.Padding() + 1
}

func (v *ViewBox) Width() int {
	return v.width
}

func (v *ViewBox) SetWidth(w int) {
	v.width = w
}

func (v *ViewBox) Height() int {
	return v.TotalLines() - 2
}

func (v *ViewBox) ContentWidth() int {
	return v.Width() - 2 - v.Padding()*2
}

func (v ViewBox) FirstLine() int {
	return v.OffsetTopStart() + v.Padding() + 1
}

func (v *ViewBox) ClearLine(row, offsetLeft, length int) {
	v.GoToCell(row, offsetLeft)
	fmt.Print(strings.Repeat(" ", length))
	v.GoToCell(row, offsetLeft)
}

func (v *ViewBox) printBottomBorder() {
	fmt.Printf("%v%s%s%s%v\n", v.Theme().Main, bottomLeftBorderCorner, strings.Repeat(horizontalBorder, v.Width()-2), bottomRightBorderCorner, "\033[0m")
}

func (v *ViewBox) printTopBorder() {
	fmt.Printf("%v%s%s%s%v\n", v.Theme().Main, topLeftBorderCorner, strings.Repeat(horizontalBorder, v.Width()-2), topRightBorderCorner, "\033[0m")
}

func (v *ViewBox) printVerticalBorder() {
	fmt.Printf("%v%s%v", v.Theme().Main, verticalBorder, "\033[0m")
}

func (v *ViewBox) GoToCell(offsetTop, offsetLeft int) {
	fmt.Printf("\033[%d;%dH", offsetTop, offsetLeft)
}

func (v *ViewBox) PrintBox() {
	v.GoToCell(v.OffsetTopStart(), v.OffsetLeftStart())

	v.printTopBorder()

	row := v.OffsetTopStart() + 1

	for row <= v.Height() {
		v.GoToCell(row, v.OffsetLeftStart())
		v.printVerticalBorder()
		v.GoToCell(row, v.Width())
		//	fmt.Printf("%s", v.backgroundColor+strings.Repeat(" ", v.ContentWidth())+"\033[0m")
		//	v.printVerticalBorder()
		v.printVerticalBorder()
		row++
	}

	v.GoToCell(row, v.OffsetLeftStart())

	v.printBottomBorder()
}

func (v *ViewBox) Notify(msg ...interface{}) {
	v.ClearLine(v.TotalLines(), 1, 200)
	v.GoToCell(v.TotalLines(), 1)
	fmt.Print(msg...)
}
