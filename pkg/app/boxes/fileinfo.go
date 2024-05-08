package boxes

import (
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type FileInfoController struct {
	*pubsub.Node
	data fs.FileInfo
	size int
	view termui.Section
}

func NewFileInfo(file fs.FileInfo, size int) *pubsub.Node {
	n := pubsub.NewNode()

	c := FileInfoController{
		n,
		file,
		size,
		termui.NewSection(),
	}
	c.view.SetBorder().SetPadding(1, 2, 1, 2).SetHeight(30).SetWidth(70).CenterVertically().CenterHorizontally()

	n.OnGlobal(pubsub.RESIZE, func() {
		c.view.CenterVertically().CenterHorizontally()
		c.render()
	})

	n.On(pubsub.RENDER, c.render)

	c.Node.OnGlobal(pubsub.T, c.render)

	n.On(pubsub.Q, func() {
		c.Node.Unlink()
		n.Passthrough(pubsub.RENDER, n.First())
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
