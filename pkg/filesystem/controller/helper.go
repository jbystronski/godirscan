package controller

import (
	"github.com/jbystronski/godirscan/pkg/common"
)

func (c *Controller) wrapInput(prompt, defaultOpt string) string {
	answ := common.WaitInput(prompt, defaultOpt, c.GoToPromptCell(), c.ErrorChan)

	defer func() {
		c.PrintBox()
		c.Alt.PrintBox()
		c.updateTotalSize()
		c.Alt.updateTotalSize()
		//	c.TotalSize(c.FsStoreAccessor.Size())
		//	c.Alt.TotalSize(c.Alt.FsStoreAccessor.Size())
	}()

	return answ
}
