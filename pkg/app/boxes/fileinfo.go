package boxes

import (
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type FileInfoController struct {
	*event.Node
	data fs.FileInfo
	size int
	view termui.Section
}

func NewFileInfo(file fs.FileInfo, size int) *event.Node {
	n := event.NewNode()

	c := FileInfoController{
		n,
		file,
		size,
		termui.NewSection(),
	}
	c.view.SetBorder().SetPadding(1, 2, 1, 2).SetHeight(30).SetWidth(70).CenterVertically().CenterHorizontally()

	n.OnGlobal(event.RESIZE, func() {
		c.view.CenterVertically().CenterHorizontally()
		c.render()
	})

	n.On(event.RENDER, c.render)

	c.Node.OnGlobal(event.T, c.render)

	n.On(event.Q, func() {
		c.Node.Unlink()
		n.Passthrough(event.RENDER, n.First())
	})

	return n
}

func (c *FileInfoController) formatRow(s string, formatting func() string) string {
	return buildString(fmtBold(formatting(), trimEnd(s, c.view.ContentWidth(), c.view.ContentWidth()-2, 2, '.')), Reset)
}

func (c *FileInfoController) render() {
	time.Sleep(time.Millisecond * 100)

	updateDimensions()

	c.view.Content = []string{
		c.formatRow("Name", ThemeAccent),
		c.formatRow(c.data.Name(), themeMain),
		c.formatRow("Size", ThemeAccent),
		c.formatRow(strings.Replace(fmtSize(c.size), " ", "", 15), themeMain),
		c.formatRow("Modified", ThemeAccent),
		c.formatRow(c.data.ModTime().String(), themeMain),
		c.formatRow("Permissions", ThemeAccent),
		c.formatRow(c.data.Mode().String(), themeMain),
	}

	c.view.Print(themeMain())
	c.view.PrintContent()
	cell(c.view.Top+c.view.Height-2, c.view.ContentStart())
	fmt.Print(c.formatRow("Press 'q' to close", ThemeAccent))
}
