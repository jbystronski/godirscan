package command_line

import (
	"fmt"
	"strings"
)

type CommandLine struct {
	output                                          []rune
	inputLine, startCol, maxVisibleOutput, position int
}

func NewCommandLine(line, startCol, maxVisibleOutput int, initialString string) *CommandLine {
	c := &CommandLine{
		inputLine:        line,
		startCol:         startCol,
		position:         startCol,
		maxVisibleOutput: maxVisibleOutput,
		output:           []rune(initialString),
	}

	c.print()

	c.GoToLineEnd()
	return c
}

func (c *CommandLine) GoToLineEnd() {
	c.position = c.getLastCol()

	c.printCursor()
}

func (c *CommandLine) GoToLineStart() {
	c.position = c.startCol

	c.printCursor()
}

func (c *CommandLine) getLastCol() int {
	return c.startCol + len(c.output)
}

func (c *CommandLine) goToCol(col int) {
	fmt.Printf("\033[%d;%dH", c.inputLine, col)
}

func (c *CommandLine) printCursor() {
	c.goToCol(c.position)
	fmt.Print("\033[?25h")
}

func (c *CommandLine) NextCol() {
	if c.position < c.getLastCol() {
		c.position++
		c.printCursor()
	}
}

func (c *CommandLine) PrevCol() {
	if c.position > c.startCol {
		c.position--
		c.printCursor()
	}
}

func (c *CommandLine) getCurrentIndex() int {
	return c.position - c.startCol
}

func (c *CommandLine) InsertChar(r rune) {
	if c.position == c.maxVisibleOutput {
		return
	}

	index := c.getCurrentIndex()

	if index > len(c.output)-1 {
		c.output = append(c.output, r)
	} else {
		c.output = append(c.output[:index], append([]rune{r}, c.output[index:]...)...)
	}

	c.position++
	c.clear()
	c.print()
}

func (c *CommandLine) DeleteChar() {
	index := c.getCurrentIndex()

	if index > len(c.output)-1 {
		return
	} else if index == len(c.output)-1 {
		c.output = c.output[:index]
	} else if index < len(c.output)-1 {
		c.output = append(c.output[:index], c.output[index+1:]...)
	}

	c.clear()
	c.print()
}

func (c *CommandLine) print() {
	c.goToCol(c.startCol)

	var visibleOutput []rune

	if c.maxVisibleOutput > len(c.output) {
		visibleOutput = c.output
	} else {
		visibleOutput = c.output[:c.maxVisibleOutput]
	}

	for _, r := range visibleOutput {
		fmt.Printf("%c", r)
	}

	c.printCursor()
}

func (c *CommandLine) Backspace() {
	index := c.getCurrentIndex()

	if index == 0 {
		return
	}

	if index > len(c.output)-1 {
		c.output = c.output[:index-1]
	} else {
		c.output = append(c.output[:index], c.output[index+1:]...)
	}
	c.position--
	c.clear()
	c.print()
}

func (c *CommandLine) clear() {
	c.goToCol(c.startCol)

	fmt.Print(strings.Repeat(" ", len(c.output)+1))
	c.goToCol(c.position)
}

func (c *CommandLine) GetOutput() string {
	var stringResult string

	for _, r := range c.output {
		stringResult += string(r)
	}
	return strings.TrimSpace(stringResult)
}
