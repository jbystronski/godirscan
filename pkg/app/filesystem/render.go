package filesystem

import (
	"fmt"

	"github.com/jbystronski/godirscan/pkg/app/data"
)

func (c *FsController) render() {
	if c.ShouldUpdateChunk() {
		c.fullRender()
	} else {
		c.updateView(c.Index())
		c.updateView(c.PrevIndex)
	}
}

func (c *FsController) fullRender() {
	clear(c.panel.Top-1, c.panel.ContentStart(), c.panel.ContentWidth())
	fmt.Print(header(trimEnd(c.root, c.panel.ContentWidth()-1, c.panel.ContentWidth()-2, 2, '.')))

	start, end := c.IndiceRange(c.Index(), c.data.Len()-1)

	c.panel.Print(themeMain())

	//	content := []string{}

	row := c.panel.ContentFirstLine()

	for i := start; i <= end; i++ {

		entry, exists := c.find(i)

		if !exists {
			continue
		}
		c.panel.PrintLine(row, c.formatEntry(entry, i))
		row++

		//	content = append(content, c.formatEntry(entry, i))
	}

	//if len(c.panel.Content) == 0 {
	//	content = append(content, c.PrintEmpty())
	//}
	// c.panel.Content = content
	// c.panel.PrintContent()

	cell(rows()-2, c.panel.ContentStart())
	// fmt.Print(fmtBold(ThemeAccent(), printSizeAsString(c.data.Size())))

	fmt.Print(totalSize(c.data.Size()))
	c.PrintBottomHelper()
}

func (c *FsController) PrintBottomHelper() {
	clear(rows(), 1, cols())
	fmt.Print(fmtBold("Press 'm' to show menu"))
}

func (c *FsController) formatEntry(entry *data.FsEntry, index int) string {
	var output string

	separator := "\u251c\u2500"

	if index == c.data.Len()-1 {
		separator = "\u2514\u2500"
	}

	t := trimEnd(fmt.Sprint(entry), c.contentWidth()-16, c.contentWidth()-16, 2, '.')

	switch true {
	case c.Index() == index && c.active:
		output = activeRow(separator, t)

	case c.selected.Exists(entry.FullPath()):
		output = selectedRow(separator, t)

	default:
		switch entry.FsType() {
		case data.DirDatatype:
			output = directory(separator, t)

		case data.SymlinkDatatype:
			output = symlink(separator, t)

		case data.FileDatatype:
			output = file(separator, t)

		}

	}

	return buildString(printSizeAsString(entry.Size()), output)
}

func (c *FsController) updateView(i int) {
	if entry, ok := c.find(i); ok {
		clear(c.GetIndexOffset(i), c.panel.ContentStart(), c.panel.ContentWidth())
		fmt.Print(c.formatEntry(entry, i))
	}
}

// func (c FsController) PrintEmpty(isPanelActive bool) {
// 	for i := p.Main.OutputFirstLine(); i <= p.Main.OutputLastLine(); i++ {
// 		clear(i, p.Main.ContentLineStart(), p.Main.ContentWidth())
// 	}

// 	cell(p.Main.OutputFirstLine(), p.Main.ContentLineStart())

// 	if isPanelActive {
// 		fmt.Print(fmtBold(themeBgHighlight(), themeHighlight(), ".."))
// 	} else {
// 		fmt.Print(fmtBold(".."))
// 	}
// }

func (c *FsController) PrintEmpty() string {
	return fmtBold(printSizeAsString(0), "\u2514\u2500", "..")
}
