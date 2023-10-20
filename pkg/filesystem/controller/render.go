package controller

import (
	"fmt"

	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func (c *Controller) fullRender() {
	c.RenderAll()
}

func (c *Controller) render() {
	if c.DataAccessor.Len() == 0 {
		c.PrintEmpty(c.Active())
		return
	}

	if c.IsChunk() {
		c.fullRender()
	} else {

		if c.PrevIndex() >= c.ChunkStart() && c.PrevIndex() <= c.ChunkEnd() {

			prev, ok := c.Find(c.PrevIndex())

			if ok {
				c.renderSingle(prev, c.PrevIndex(), false)
			}

		}

		curr, ok := c.Find(c.Index())

		if ok {
			c.renderSingle(curr, c.Index(), true)
		}

	}
}

func (c *Controller) renderSingle(e *filesystem.FsEntry, index int, isIndexActive bool) {
	c.ClearLine(c.Line(index), c.ContentLineStart(), c.ContentWidth())
	var output string
	sep := c.Separator(index == c.DataAccessor.Len()-1)
	t := c.TrimEnd(fmt.Sprint(e), c.ContentWidth()-16, c.ContentWidth()-16, 2, '.')

	switch true {

	case c.Active() && isIndexActive:
		output = c.ActiveRow(sep, t)

	case c.selected.Exists(e.FullPath()):
		output = c.SelectedRow(sep, t)

	default:
		switch e.FsType() {
		case filesystem.Dir:
			output = c.Directory(sep, t)

		case filesystem.Symlink:
			output = c.Symlink(sep, t)

		case filesystem.File:
			output = c.File(sep, t)

		default:
			output = c.File(sep, t)
		}
	}

	fmt.Print(c.PrintSizeAsString(e.Size()), output)
}

func (c *Controller) RenderAll() {
	c.ClearLine(1, c.ContentLineStart()-1, c.Width())

	header := c.TrimEnd(c.path, c.Width()-1, c.Width()-1, 2, '.')

	fmt.Print(c.Header(header))

	if c.DataAccessor.Len() == 0 {
		c.PrintEmpty(c.Active())
		return
	}

	c.GoToCell(c.Line(c.Index()), c.ContentLineStart())

	line := c.FirstLine()

	for i := c.ChunkStart(); i <= c.ChunkEnd(); i++ {
		file, ok := c.Find(i)

		if !ok {
			continue
		}

		c.renderSingle(file, i, i == c.Index())

		line++
		c.GoToCell(line, c.ContentLineStart())
	}

	for line <= c.OutputLastLine() {
		c.ClearLine(line, c.ContentLineStart(), c.ContentWidth())
		line++
	}
	c.GoToCell(c.TotalLines()-2, c.ContentLineStart())

	fmt.Print(c.TotalSize(c.DataAccessor.Size()))

	c.GoToCell(c.TotalLines(), 1)
}

func (c *Controller) restoreScreen() {
	common.ClearScreen()
	c.PrintBox()
	c.fullRender()
	c.Alt.PrintBox()
	c.Alt.fullRender()
}
