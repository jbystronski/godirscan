package controller

import (
	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func (c *Controller) right() {
	if c.DataAccessor.Len() == 0 {
		return
	}

	en, ok := c.Find(c.Index())

	if !ok {
		return
	}

	if en.FsType() == filesystem.Dir {
		c.backtrace.Set(c.path, c.Index())
		c.ChunkNavigator.Reset()
		c.SetStore(en.FullPath())
	}

	//	c.SetInitialData(dir)

	//	c.SetStore(c.Data().Find(c.Index()).FullPath())

	//	c.DefaultSort()
	//	c.Reset()
	// c.SetChunk(c.Data().Len())
	// c.fullRender()
}
