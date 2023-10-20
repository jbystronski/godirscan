package controller

import (
	"fmt"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
)

type keyController struct {
	key,
	description string
}

var keys = []keyController{
	{"enter", "Enter directory or execute file with deault application"},
	{"ctrl e", "Execute command inside the current directory"},
	{"delete / F8", "Delete selected entries"},
	{"ctrl r", "Rename file / directory"},
	{"F7", "Create new directory"},
	{"ctrl w", "Create new file"},
	{"ctrl f", "Search"},
	{"ctrl s", "Sort entries within directory"},
	{"F6", "Move selected into the current directory"},
	{"ctrl v", "Copy selected into the current directory"},
	{"F9", "Scan new directory"},
	{"home / end", "Go to first / last entry"},
	{"pg down / pg up", "Jump to the top / end of visible entries"},
	{"left / \033[1C", "Go to parent or child directory"},
	{"up / down", "Navigate up and down"},
	{"insert", "Select entry"},
	{"ctrl a", "Select / deselect all entries"},
	{"ctrl g", ""},
	{"ctrl 4", "Edit file"},
	{"esc", "Stop currently running process (copying, searching, ...)"},
	{"ctrl /", "Change theme"},
	{"ctrl 5", "Print help"},
	{"F10, ctrl c / q", "Quit"},
}

func (c Controller) printHelpMenu(startIndex, endIndex int) {
	row := 11

	keysWidth := 20
	descriptionsWidth := 71

	for i := startIndex; i <= endIndex; i++ {

		c.GoToCell(row, 1)
		name := c.BuildString(c.Theme().Accent, "\033[1m", c.AlignRight(keysWidth, keys[i].key, " "), "\033[0m")

		fmt.Print(name)
		c.GoToCell(row, keysWidth+1)

		description := c.BuildString("\033[1m", c.AlignRight(descriptionsWidth, keys[i].description, " "), "\033[0m")

		fmt.Print(description)
		row++
	}
}

func (c *Controller) printHelp() {
	// resizeChan := make(chan os.Signal, 1)

	// signal.Notify(resizeChan, syscall.SIGWINCH)

	common.ClearScreen()
	c.PrintBanner()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-resizeChan:
	// 			c.restoreScreen()
	// 			return
	// 		}
	// 	}
	// }()

	startIndex := 0
	endIndex := 9

	common.ClearScreen()
	c.PrintBanner()

	c.printHelpMenu(startIndex, endIndex)

	for {
		_, key, err := k.GetKey()
		if err != nil {
			c.ErrorChan <- err
		}

		switch key {
		case k.KeyEsc:
			c.restoreScreen()
			return
		case k.KeyArrowDown:
			if endIndex < len(keys)-1 {
				startIndex++
				endIndex++
				c.printHelpMenu(startIndex, endIndex)
			}
		case k.KeyArrowUp:
			if startIndex != 0 {
				startIndex--
				endIndex--
				c.printHelpMenu(startIndex, endIndex)
			}
		}

	}
}
