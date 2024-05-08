package filesystem

import (
	"strings"

	"github.com/eiannone/keyboard"
)

func (c *FsController) goTo() {
	if c.data.Len() == 0 {
		return
	}

	var matches []int
	var curr int

	test := c.getInput("Go to", "")

	if test == "" {
		return
	}

	for index, en := range c.data.All() {
		if strings.Contains(en.Name(), test) {
			matches = append(matches, index)
		}
	}

	if len(matches) == 0 {
		return
	}

	c.Jump(c.PrevIndex, matches[curr], c.data.Len())

	c.render()

	msg := "Press up and down to jump between found entries, or esc to resume navigation"

	openCommandLine()
	//	termui.NewTerminal().CommandLineOpen = true
	clearPrompt()
	printInfo(msg)

	for {
		select {
		case ev := <-receiveCommandLine():
			switch true {
			case ev.Key != 0:
				switch ev.Key {
				case keyboard.KeyArrowRight:
					c.execute()
				case keyboard.KeyArrowUp, keyboard.KeyArrowLeft:
					if curr == 0 {
						continue
					}
					prev := matches[curr]
					curr--
					c.Jump(prev, matches[curr], c.data.Len())

					c.render()
					clearPrompt()
					printInfo(msg)

				case keyboard.KeyArrowDown, keyboard.KeyInsert:
					if curr == len(matches)-1 {
						continue
					}

					prev := matches[curr]

					if ev.Key == keyboard.KeyInsert {
						en, ok := c.data.Find(prev)

						if !ok {
							return
						}

						if c.root == c.selected.BasePath() {
							c.selected.Clear()
						}
						c.selected.Toggle(en.FullPath())
					}

					curr++
					c.Jump(prev, matches[curr], c.data.Len())

					c.render()
					clearPrompt()
					printInfo(msg)

				case keyboard.KeyEsc:
					closeCommandLine()
					clearPrompt()
					c.render()
					return

				}
			}
		}
	}
}
