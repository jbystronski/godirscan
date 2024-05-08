package menu

import (
	"fmt"
	"time"

	"github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type Dimensions struct {
	PaddingLeft   int
	PaddingRight  int
	PaddingTop    int
	PaddingBottom int
	Width         int
	Height        int
}

type MenuController struct {
	*pubsub.Node
	*pubsub.Broker
	termui.Navigator
	Options []MenuOption
	View    termui.Section
	index,
	startIndex,
	endIndex,
	CurrentLine int
}

func NewMenuController(options []MenuOption, dimensions Dimensions) *MenuController {
	c := MenuController{
		pubsub.NewNode(),
		pubsub.SingleBroker(),
		termui.Navigator{},
		options,
		termui.NewSection(),

		0,
		0,
		0,
		0,
	}

	c.Node.On(pubsub.Q, func() {
		//	termui.ClearScreen()
		c.Node.Unlink()
		c.Node.Passthrough(pubsub.RENDER, c.Node.Prev)
	})

	c.Node.On(pubsub.ARROW_DOWN, func() {
		if c.NextEntry() {
			c.render()
		}
	})

	c.Node.On(pubsub.ARROW_UP, func() {
		if c.PrevEntry() {
			c.render()
		}
	})

	c.Node.On(pubsub.PG_UP, func() {
		if len(c.Options) > 0 && c.MovePgUp() {
			c.render()
		}
	})

	c.Node.On(pubsub.PG_DOWN, func() {
		if c.MovePgDown(c.View.ContentLines()) {
			c.render()
		}
	})

	c.Node.On(pubsub.HOME, func() {
		if c.FirstEntry() {
			c.render()
		}
	})

	c.Node.On(pubsub.END, func() {
		if c.LastEntry() {
			c.render()
		}
	})

	c.Node.OnGlobal(pubsub.RESIZE, func() {
		c.View.CenterVertically().CenterHorizontally()

		c.Navigator.MinOffset = c.View.OutputFirstLine()
		c.Navigator.MaxOffset = c.View.OutputLastLine() - 1
		c.print()
	})

	c.Node.OnGlobal(pubsub.T, func() {
		time.Sleep(time.Millisecond * 100)

		t := termui.NewTerminal()
		t.UpdateDimensions()

		c.fullRender()
	})

	c.Node.On(pubsub.RENDER, c.print)

	c.View.SetBorder().SetPadding(1, 1, 1, 2).SetHeight(dimensions.Height).SetWidth(dimensions.Width).CenterVertically().CenterHorizontally()
	c.Navigator.MinOffset = c.View.OutputFirstLine()
	c.Navigator.MaxOffset = c.View.OutputLastLine() - 1
	c.Navigator.TotalEntries = len(options)

	return &c
}

func (m *MenuController) RunDefault(e pubsub.Event) {
	var target *pubsub.Node
	if m.Prev.HasLocal(e) {
		target = m.Prev
	} else {
		target = m.First()
	}

	m.Unlink()
	m.Passthrough(pubsub.RENDER, target)
	m.Passthrough(e, target)
}

func (m *MenuController) print() {
	time.Sleep(time.Millisecond * 100)

	t := termui.NewTerminal()
	t.UpdateDimensions()

	m.render()
}

func (m *MenuController) fullRender() {
	start, end := m.IndiceRange(m.Index(), len(m.Options)-1)

	m.View.Print(global.ThemeMain())

	content := []string{}

	for i := start; i <= end; i++ {
		content = append(content, global.FmtBold(global.ThemeAccent(), m.Options[i].Label, global.ThemeMain(), " ", m.Options[i].Description))
	}

	m.View.Content = content
	m.View.PrintContent()

	m.formatRow(m.GetIndexOffset(m.Index()), m.Index(), m.highlight)
}

func (m *MenuController) update() {
	m.formatRow(m.GetIndexOffset(m.Index()), m.Index(), m.highlight)
	m.formatRow(m.GetIndexOffset(m.PrevIndex), m.PrevIndex, m.withoutHighligth)
}

func (m *MenuController) render() {
	if m.ShouldUpdateChunk() {
		m.fullRender()
	} else {
		m.update()
	}
}

func (m *MenuController) withoutHighligth(opt MenuOption) {
	fmt.Print(global.FmtBold(global.ThemeAccent(), opt.Label, global.ThemeMain(), " ", opt.Description))
}

func (m *MenuController) highlight(opt MenuOption) {
	fmt.Print(global.FmtBold(global.ThemeBgHighlight(), global.ThemeHighlight(), opt.Label, " ", opt.Description))
}

func (m *MenuController) formatRow(offsetX, index int, format func(opt MenuOption)) {
	global.Clear(offsetX, m.View.ContentStart(), m.View.ContentWidth())
	format(m.Options[index])
}
