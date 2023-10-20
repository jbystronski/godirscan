package common

import (
	"context"
	"sync"
)

type CancelContext struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

func (c *CancelContext) Create() {
	c.Ctx, c.CancelFunc = context.WithCancel(context.Background())
}

func (c *CancelContext) Cancel() {
	c.CancelFunc()
}

func (c *CancelContext) Observe(task func()) {
	var wg sync.WaitGroup
	done := make(chan struct{}, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {

			case <-done:

				return

			case <-c.Ctx.Done():

				return
			//	done <- struct{}{}
			default:

				task()
				done <- struct{}{}

			}
		}
	}()

	wg.Wait()
}
