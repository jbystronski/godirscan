package controller

import (
	"errors"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
)

func (c *Controller) goTo() {
	if c.DataAccessor.Len() > 0 {
		var matches []int
		var curr int

		test := c.wrapInput("Go to: ", "")

		if test == "" {
			return
		}

		for index, en := range c.DataAccessor.All() {
			if strings.Contains(en.Name(), test) {
				matches = append(matches, index)
			}
		}

		if len(matches) > 0 {

			c.JumpTo(c.PrevIndex(), matches[0], c.DataAccessor.Len())

			c.render()
			c.GoToCell(c.TotalLines()-2, c.ContentLineStart())

		} else {
			c.ErrorChan <- errors.New("no matches found")

			return
		}

	MatchLoop:
		for {

			_, key, err := keyboard.GetKey()
			if err != nil {
				c.ErrorChan <- err
				break
			}
			switch key {

			case keyboard.KeyInsert:

				en, ok := c.Find(c.Index())

				if !ok {
					return
				}

				if c.path == c.selected.BasePath() {
					c.selected.Clear()
				}
				c.selected.Toggle(en.FullPath())
				c.render()

			case keyboard.KeyArrowUp, keyboard.KeyArrowLeft:
				if curr == 0 {
					continue
				}
				prev := matches[curr]
				curr--
				c.JumpTo(prev, matches[curr], c.DataAccessor.Len())

				c.render()
				common.Cell(c.TotalLines()-2, c.ContentLineStart())

			case keyboard.KeyArrowDown, keyboard.KeyArrowRight:
				if curr == len(matches)-1 {
					continue
				}

				prev := matches[curr]
				curr++
				c.JumpTo(prev, matches[curr], c.DataAccessor.Len())

				c.render()
				common.Cell(c.TotalLines()-2, c.ContentLineStart())

			case keyboard.KeyEsc, keyboard.KeyEnter:
				common.ClearLine()
				common.CarriageReturn()
				break MatchLoop
			}

		}

		c.render()
		c.Alt.render()

	}
}
