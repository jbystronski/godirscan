package app

import (
	"os"
	"sync"

	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/app/filesystem"
	"github.com/jbystronski/godirscan/pkg/app/menu"
	"github.com/jbystronski/godirscan/pkg/global/event"
	"github.com/jbystronski/pubsub"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

var (
	startOnce sync.Once
	startNode *pubsub.Node
)

func NewStart() *pubsub.Node {
	startOnce.Do(func() {
		n := pubsub.NewNode(pubsub.GlobalBroker())

		n.On(event.RENDER, printStartScreen)

		n.On(event.S, func() {
			printStartPrompt(n)
		})

		n.OnGlobal(event.T, printStartScreen)

		n.OnGlobal(event.RESIZE, printStartScreen)

		n.On(event.ESC, func() {
			m := menu.QuitMenu()
			// m.Watch()
			n.LinkTo(m.Node)
			n.Passthrough(event.RENDER, n.Next())
		})

		startNode = n
	})

	return startNode
}

func printStartScreen() {
	cls()

	banner := banner()
	s := termui.NewSection()

	s.SetHeight(len(banner)).SetWidth(86).CenterHorizontally().CenterVertically()
	s.Content = banner
	x := s.PrintContent()

	line := "v1.0 by pogodisco"

	s.Width = strlen(line)
	s.Top = x
	s.CenterHorizontally()
	s.Content = []string{fmtBold(ThemeAccent(), line, termui.Reset)}
	x = s.PrintContent()

	line = "Press 's' to scan a directory, 'esc' to quit"

	s.Width = strlen(line)
	s.Top = x
	s.CenterHorizontally()
	s.Content = []string{fmtBold(ThemeAccent(), line)}
	s.PrintContent()

	hideCursor()
}

func printStartPrompt(n *pubsub.Node) {
	clear(rows(), 1, cols())

	cmd := termui.NewCommandLine(rows(), 1, "Scan a directory", fmtPrompt, config.Running().DefaultRootDirectory)
	root := cmd.WaitInput()

	if root == "" {
		n.Passthrough(event.RENDER, n)
	} else {

		if _, err := os.Stat(root); err != nil {
			n.Publish("err", pubsub.Message(err.Error()))
			return
		}

		fs := filesystem.New(root)

		// fs.Watch()
		n.Prev().LinkTo(fs)
		n.Prev().Passthrough(event.INIT, n.Prev().Next())

	}
}
