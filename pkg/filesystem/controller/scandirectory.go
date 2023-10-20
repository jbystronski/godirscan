package controller

import "github.com/jbystronski/godirscan/pkg/common"

func (c *Controller) scanDir() {
	dir := c.wrapInput("Scan directory: ", common.Cfg.DefaultRootDirectory)

	if dir == "" {
		return
	}

	c.Reset()
	c.backtrace.Clear()
	c.cache.Clear()

	c.SetStore(dir)

	c.Alt.SetStore(c.Alt.path)
}
