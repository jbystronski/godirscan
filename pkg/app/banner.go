package app

import (
	"fmt"
	"strings"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type Banner struct {
	Height, Width, x, y int
}

func (b Banner) Print(offsetX, offsetY int) {
	b.x = offsetX
	b.y = offsetY

	cell(b.x, b.y)

	printRow := func(parts ...string) {
		fmt.Print(trimEnd(buildString(parts...), cols(), rows(), 0, ' '))
		b.x += 1
		cell(b.x, b.y)
	}

	slash := termui.Slash
	block := termui.Block
	segment := termui.Segment
	line := termui.Line
	space := termui.Space
	bslash := termui.Backslash
	repeat := strings.Repeat

	printRow(themeMain(), repeat(" ", 4), slash, repeat(block, 8), slash, repeat(space, 74))

	printRow(repeat(" ", 3), segment, repeat(line, 6), repeat(space, 2), slash, repeat(block, 7), slash, space, slash, repeat(block, 6), bslash, repeat(space, 2), slash, block, slash, space, slash, repeat(block, 6), bslash, repeat(space, 2), slash, repeat(block, 6), slash, space, slash, repeat(block, 6), slash, space, slash, repeat(block, 7), slash, space, slash, repeat(block, 2), bslash, repeat(space, 2), slash, block, slash, repeat(space, 3))

	printRow(repeat(" ", 2), segment, repeat(space, 2), slash, repeat(block, 3), slash, space, segment, repeat(line, 3), segment, space, segment, repeat(line, 3), segment, space, segment, space, segment, repeat(line, 3), segment, space, slash, repeat(block, 5), bslash, repeat(space, 2), segment, repeat(line, 5), space, segment, repeat(line, 3), segment, space, segment, bslash, block, bslash, segment, repeat(space, 4))

	printRow(repeat(" ", 1), segment, repeat(space, 3), line, segment, space, segment, repeat(space, 3), segment, space, segment, repeat(space, 3), segment, space, segment, space, slash, repeat(block, 7), repeat(line, 0), repeat(space, 0), slash, repeat(space, 2), repeat(line, 4), slash, block, slash, space, segment, repeat(space, 6), slash, repeat(block, 7), slash, space, segment, repeat(space, 2), bslash, repeat(block, 2), slash, repeat(space, 5))

	printRow(repeat(" ", 0), slash, repeat(block, 8), slash, space, slash, repeat(block, 7), slash, space, slash, repeat(block, 7), slash, space, slash, block, slash, space, segment, repeat(line, 4), bslash, block, bslash, space, slash, repeat(block, 5), slash, space, slash, repeat(block, 6), slash, space, segment, repeat(line, 3), segment, space, segment, repeat(space, 3), segment, repeat(space, 6))

	printRow(repeat(" ", 0), repeat(line, 9), repeat(space, 2), repeat(line, 8), repeat(space, 2), repeat(line, 8), repeat(space, 2), line, line, repeat(space, 2), repeat(line, 2), repeat(space, 6), repeat(line, 2), space, repeat(line, 6), repeat(space, 2), repeat(line, 7), repeat(space, 2), repeat(line, 2), repeat(space, 4), repeat(line, 2), repeat(space, 2), repeat(line, 2), repeat(space, 4), repeat(line, 2), repeat(space, 7), termui.Reset)
}

func banner() []string {
	slash := termui.Slash
	block := termui.Block
	segment := termui.Segment
	line := termui.Line
	space := termui.Space
	bslash := termui.Backslash
	repeat := strings.Repeat

	content := []string{}

	content = append(content, buildString(themeMain(), repeat(" ", 4), slash, repeat(block, 8), slash, repeat(space, 74)))
	content = append(content, buildString(repeat(" ", 3), segment, repeat(line, 6), repeat(space, 2), slash, repeat(block, 7), slash, space, slash, repeat(block, 6), bslash, repeat(space, 2), slash, block, slash, space, slash, repeat(block, 6), bslash, repeat(space, 2), slash, repeat(block, 6), slash, space, slash, repeat(block, 6), slash, space, slash, repeat(block, 7), slash, space, slash, repeat(block, 2), bslash, repeat(space, 2), slash, block, slash, repeat(space, 3)))
	content = append(content, buildString(repeat(" ", 2), segment, repeat(space, 2), slash, repeat(block, 3), slash, space, segment, repeat(line, 3), segment, space, segment, repeat(line, 3), segment, space, segment, space, segment, repeat(line, 3), segment, space, slash, repeat(block, 5), bslash, repeat(space, 2), segment, repeat(line, 5), space, segment, repeat(line, 3), segment, space, segment, bslash, block, bslash, segment, repeat(space, 4)))

	content = append(content, buildString(repeat(" ", 1), segment, repeat(space, 3), line, segment, space, segment, repeat(space, 3), segment, space, segment, repeat(space, 3), segment, space, segment, space, slash, repeat(block, 7), repeat(line, 0), repeat(space, 0), slash, repeat(space, 2), repeat(line, 4), slash, block, slash, space, segment, repeat(space, 6), slash, repeat(block, 7), slash, space, segment, repeat(space, 2), bslash, repeat(block, 2), slash, repeat(space, 5)))

	content = append(content, buildString(repeat(" ", 0), slash, repeat(block, 8), slash, space, slash, repeat(block, 7), slash, space, slash, repeat(block, 7), slash, space, slash, block, slash, space, segment, repeat(line, 4), bslash, block, bslash, space, slash, repeat(block, 5), slash, space, slash, repeat(block, 6), slash, space, segment, repeat(line, 3), segment, space, segment, repeat(space, 3), segment, repeat(space, 6)))

	content = append(content, buildString(repeat(" ", 0), repeat(line, 9), repeat(space, 2), repeat(line, 8), repeat(space, 2), repeat(line, 8), repeat(space, 2), line, line, repeat(space, 2), repeat(line, 2), repeat(space, 6), repeat(line, 2), space, repeat(line, 6), repeat(space, 2), repeat(line, 7), repeat(space, 2), repeat(line, 2), repeat(space, 4), repeat(line, 2), repeat(space, 2), repeat(line, 2), repeat(space, 4), repeat(line, 2), repeat(space, 7), termui.Reset))

	return content
}
