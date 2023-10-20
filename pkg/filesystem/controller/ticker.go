package controller

import (
	"sync"
	"time"
)

func (c *Controller) ObserveTicker(ctxDone <-chan struct{}, interval time.Duration, intervalCallback func()) (*sync.WaitGroup, chan struct{}) {
	var ticker time.Ticker
	var init bool
	var wg sync.WaitGroup
	done := make(chan struct{})
	ticker = *time.NewTicker(interval)
	init = true

	// c.Tickable.Start(interval)

	// wg.Add(1)

	go func() {
		defer func() {
			ticker.Stop()
			//	c.Tickable.Stop()

			//	wg.Done()
		}()
		for {
			select {

			case err := <-c.internalErrChan:

				c.ErrorChan <- err
				return

			case <-ticker.C:
				if init {
					intervalCallback()
				}

			case <-ctxDone:
				return

			case <-done:

				return
			}
		}
	}()

	return &wg, done
}
