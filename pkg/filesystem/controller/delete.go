package controller

import (
	"errors"
	"os"
	"time"
)

func (c *Controller) deleteEntries() {
}

func (c *Controller) delete() {
	answ := c.wrapInput("Delete selected entries", "y")

	if answ != "y" {
		// c.render()
		// c.Alt.render()
		return
	}

	// doneChan := make(chan struct{}, 1)
	// messageChan := make(chan string)

	//go func() {
	//	common.PrintProgress(doneChan, messageChan, 1, common.NumVisibleLines())
	//}()

	if c.path == c.selected.BasePath() {
		c.selected.Clear()
	}
	// var wg sync.WaitGroup

	//	_, tickerDone := c.ObserveTicker(c.IsCancelled(), time.Millisecond*200, c.fullRender)

	for key := range c.selected.Self() {
		//	wg.Add(1)
		// go func(key string) {
		//	defer wg.Done()

		// c.WithContext(func() {
		//	defer func() {
		// messageChan <- fmt.Sprint("Deleted ", key, " ")
		//		}()

		if error := os.RemoveAll(key); error != nil {
			if !errors.Is(error, os.ErrNotExist) {
				c.ErrorChan <- error
				return
			}
		}
		//	})
		//	}(key)
	}
	//	wg.Wait()

	// tickerDone <- struct{}{}
	c.selected.Clear()
	c.cache.Clear()
	//	doneChan <- struct{}{}

	c.SetStore(c.path)

	if c.Alt.path == c.path {
		time.Sleep(time.Millisecond * 200)
		c.Alt.SetStore(c.path)
	}

	//	c.Alt.SetStore(c.Alt.FsStoreAccessor.Name())

	// if c.Data().Len() == 0 {
	// 	c.Reset()
	// } else if c.Index() > c.Data().Len()-1 {
	// 	c.SetIndex(c.Data().Len() - 1)
	// }

	//	c.fullRender()

	// c.fullRender()
	// c.Alt.fullRender()
}
