package viewbox

import (
	"fmt"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
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
	row,
	col,
	offsetTopStart,
	offsetLeftStart,
	offsetTop,
	offsetBottom,
	totalLines,
	width,
	padding int
	theme           *common.Theme
	backgroundColor string
}

func (v *ViewBox) Position() (int, int) {
	return v.col, v.row
}

func (v *ViewBox) Theme() *common.Theme {
	return v.theme
}

func (v *ViewBox) SetTheme(t *common.Theme) {
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

func (v ViewBox) BuildString(substrings ...string) string {
	var builder strings.Builder

	for _, v := range substrings {
		builder.WriteString(v)
	}

	return builder.String()
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
	s := v.BuildString(v.Theme().Main, bottomLeftBorderCorner, strings.Repeat(horizontalBorder, v.Width()-2), bottomRightBorderCorner, "\033[0m", "\n")

	fmt.Print(s)

	// fmt.Printf("%v%s%s%s%v\n", v.Theme().Main, bottomLeftBorderCorner, strings.Repeat(horizontalBorder, v.Width()-2), bottomRightBorderCorner, "\033[0m")
}

func (v *ViewBox) printTopBorder() {
	s := v.BuildString(v.Theme().Main, topLeftBorderCorner, strings.Repeat(horizontalBorder, v.Width()-2), topRightBorderCorner, "\033[0m", "\n")

	fmt.Print(s)

	// fmt.Printf("%v%s%s%s%v\n", v.Theme().Main, topLeftBorderCorner, strings.Repeat(horizontalBorder, v.Width()-2), topRightBorderCorner, "\033[0m")
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
		v.GoToCell(row, v.offsetLeftStart+v.Width()-1)

		v.printVerticalBorder()
		row++
	}

	v.GoToCell(row, v.OffsetLeftStart())

	v.printBottomBorder()
}

func (v ViewBox) AlignLeft(maxWidth int, st, padding string) string {
	len, _ := v.splitString(st)

	if maxWidth > len {
		return v.BuildString(st, strings.Repeat(padding, maxWidth-len))
	}
	return st
}

func (v ViewBox) AlignRight(maxWidth int, st, padding string) string {
	len, _ := v.splitString(st)

	if maxWidth > len {
		return v.BuildString(strings.Repeat(padding, maxWidth-len), st)
	}
	return st
}

func (v ViewBox) splitString(s string) (int, []rune) {
	slice := []rune(s)

	return len(slice), slice
}

func (v ViewBox) breakString(st string, maxWidth int) []string {
	l, runeSlice := v.splitString(st)
	if maxWidth > l {
		return []string{st}
	}

	parts := []string{}

	start := 0
	end := maxWidth

	for {
		if start >= l {
			break
		}

		if start+maxWidth > l {
			parts = append(parts, string(runeSlice[start:]))
			break
		}

		parts = append(parts, string(runeSlice[start:end]))
		start = end
		end = end + maxWidth

	}

	return parts
}

func (v ViewBox) TrimEnd(input string, maxLen, trimLen, swaps int, swapChar rune) string {
	toTrim := []rune(input)

	if len(toTrim) > maxLen {

		toTrim = toTrim[0:trimLen]
		toTrim = v.SwapChars(toTrim, swaps, swapChar)

	}

	return string(toTrim)
}

// func (v ViewBox) TrimFrom(input string, maxLen, trimLen, swaps int, swapChar rune) string {
// 	inputLen, inputSlice := v.splitString(input)

// 	if inputLen > maxLen {

// 	}

// }

func (v ViewBox) SwapChars(chars []rune, swaps int, swapChar rune) []rune {
	for swaps > 0 {

		chars[len(chars)-swaps] = swapChar

		swaps--
	}

	return chars
}

// func (v ViewBox) MinifyPath(path string, maxLen int) string {
// 	len, chars := v.splitString(path)

// 	if len <= maxLen {
// 		return path
// 	}
// }

func (v ViewBox) PrintBanner() {
	const (
		doubleHorizontal = "\u2550"
		slash            = "\u2571"
		doubleVertical   = "\u2551"
		space            = " "
		reset            = "\033[0m"
		Line             = "\u2594"
		texture          = "\u2591"
		backslash        = "\u2572"
	)

	segment := strings.Join([]string{slash, texture, slash}, "")

	helper := func(parts ...string) {
		fmt.Print(v.TrimEnd(v.BuildString(parts...), common.NumVisibleCols(), common.NumVisibleCols(), 0, ' '))
	}

	fmt.Print(v.theme.Main)

	helper("\u2554", strings.Repeat(doubleHorizontal, 90), "\u2557\n")
	helper(doubleVertical, strings.Repeat(space, 90), doubleVertical, "\n")

	helper(doubleVertical, strings.Repeat(" ", 6), slash, strings.Repeat(texture, 8), slash, strings.Repeat(space, 74), doubleVertical, "\n")

	helper(doubleVertical, strings.Repeat(" ", 5), segment, strings.Repeat(Line, 6), strings.Repeat(space, 2), slash, strings.Repeat(texture, 7), slash, strings.Repeat(space, 1), slash, strings.Repeat(texture, 6), backslash, strings.Repeat(space, 2), slash, strings.Repeat(texture, 1), slash, strings.Repeat(space, 1), slash, strings.Repeat(texture, 6), backslash, strings.Repeat(space, 2), slash, strings.Repeat(texture, 6), slash, strings.Repeat(space, 1), slash, strings.Repeat(texture, 6), slash, strings.Repeat(space, 1), slash, strings.Repeat(texture, 7), slash, strings.Repeat(space, 1), slash, strings.Repeat(texture, 2), backslash, strings.Repeat(space, 2), slash, strings.Repeat(texture, 1), slash, strings.Repeat(space, 3), doubleVertical, "\n")

	helper(doubleVertical, strings.Repeat(" ", 4), segment, strings.Repeat(space, 2), slash, strings.Repeat(texture, 3), slash, strings.Repeat(space, 1), segment, strings.Repeat(Line, 3), segment, space, segment, strings.Repeat(Line, 3), segment, strings.Repeat(space, 1), segment, space, segment, strings.Repeat(Line, 3), segment, strings.Repeat(space, 1), slash, strings.Repeat(texture, 5), backslash, strings.Repeat(space, 2), segment, strings.Repeat(Line, 5), strings.Repeat(space, 1), segment, strings.Repeat(Line, 3), segment, strings.Repeat(space, 1), segment, backslash, texture, backslash, segment, strings.Repeat(space, 4), doubleVertical, "\n")

	helper(doubleVertical, strings.Repeat(" ", 3), segment, strings.Repeat(space, 3), strings.Repeat(Line, 1), segment, space, segment, strings.Repeat(space, 3), segment, space, segment, strings.Repeat(space, 3), segment, strings.Repeat(space, 1), segment, space, slash, strings.Repeat(texture, 7), strings.Repeat(Line, 0), strings.Repeat(space, 0), slash, strings.Repeat(space, 2), strings.Repeat(Line, 4), slash, texture, slash, space, segment, strings.Repeat(space, 6), slash, strings.Repeat(texture, 7), slash, space, segment, strings.Repeat(space, 2), backslash, strings.Repeat(texture, 2), slash, strings.Repeat(space, 5), doubleVertical, "\n")

	helper(doubleVertical, strings.Repeat(" ", 2), slash, strings.Repeat(texture, 8), slash, space, slash, strings.Repeat(texture, 7), slash, space, slash, strings.Repeat(texture, 7), slash, space, slash, texture, slash, strings.Repeat(space, 1), segment, strings.Repeat(Line, 4), backslash, texture, backslash, strings.Repeat(space, 1), slash, strings.Repeat(texture, 5), slash, space, slash, strings.Repeat(texture, 6), slash, space, segment, strings.Repeat(Line, 3), segment, space, segment, strings.Repeat(space, 3), segment, strings.Repeat(space, 6), doubleVertical, "\n")

	helper(doubleVertical, strings.Repeat(" ", 2), strings.Repeat(Line, 9), strings.Repeat(space, 2), strings.Repeat(Line, 8), strings.Repeat(space, 2), strings.Repeat(Line, 8), strings.Repeat(space, 2), Line, Line, strings.Repeat(space, 2), strings.Repeat(Line, 2), strings.Repeat(space, 6), strings.Repeat(Line, 2), space, strings.Repeat(Line, 6), strings.Repeat(space, 2), strings.Repeat(Line, 7), strings.Repeat(space, 2), strings.Repeat(Line, 2), strings.Repeat(space, 4), strings.Repeat(Line, 2), strings.Repeat(space, 2), strings.Repeat(Line, 2), strings.Repeat(space, 4), strings.Repeat(Line, 2), strings.Repeat(space, 7), doubleVertical, "\n")

	helper("\u255A", strings.Repeat(doubleHorizontal, 90), "\u255D")

	fmt.Print(reset)

	fmt.Print("\n\n")
}
