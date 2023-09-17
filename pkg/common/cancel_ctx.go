package common

import "context"

type CancelCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
	init   bool
}

func (c *CancelCtx) Create() {
	if !c.init {
		c.ctx, c.cancel = context.WithCancel(context.Background())
		c.init = true
	}
}

func (c *CancelCtx) Cancel() {
	if c.init {
		c.cancel()
		c.init = false
	}
}

func (c *CancelCtx) IsCancelled() <-chan struct{} {
	return c.ctx.Done()
}
