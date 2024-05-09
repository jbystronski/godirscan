package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jbystronski/godirscan/pkg/app/boxes"
	"github.com/jbystronski/godirscan/pkg/app/data"
	"github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/global/event"

	"github.com/jbystronski/pubsub"
)

func (c *FsController) scanDir(defaultDir string) {
	dir := c.getInput("Scan directory", defaultDir)

	if dir == "" {
		return
	}

	if _, err := os.Stat(dir); err != nil {

		c.sendError(err)

		return

	}

	c.Navigator.Reset()
	c.backtrace.Clear()
	c.cache.Clear()
	//	c.ui.Main.Print()
	c.setStore(dir)

	// c.alt.setStore(c.alt.root)
}

func (c *FsController) setStore(path string) error {
	// cls()

	err := c.setInitialData(path)
	if err != nil {
		return err
	}

	// progress := progress.New(c.ctx.CancelFunc)
	// progress.Watch()
	// c.LinkTo(progress)
	// c.Passthrough(e.RENDER, c.Next)

	c.calculateStoreSize()

	// c.Next.Unlink()

	time.Sleep(time.Millisecond * 130)

	c.sort()
	// c.alt.sort()
	c.TotalEntries = c.data.Len()
	c.fullRender()

	// c.alt.fullRender()
	// c.restoreScreen()
	// c.restartView()

	return nil
}

func (c *FsController) setInitialData(path string) error {
	global.ResolveUserDirectory(&path)

	c.root = path

	for _, v := range c.data.All() {
		v.SetSize(0)
		v.SetName("")
		c.pool.Put(v)
	}

	dc, err := global.ReadIgnorePermission(path)
	if err != nil {
		return err
	}

	c.data.Reset()

	for _, en := range dc {
		// entry := c.pool.Get()

		info, err := os.Lstat(filepath.Join(path, en.Name()))
		if err != nil {
			continue
		}

		c.insertEntry(path, info)

	}

	if c.data.Len() > 0 && c.Index() > c.data.Len()-1 {
		c.SetIndex(c.data.Len() - 1)
	}

	// c.restartView()

	return nil
}

func (c *FsController) calculateStoreSize() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// c.ctx.Create()

		ctx := global.NewCancelContext()

		progress := boxes.NewProgressBox(ctx.CancelFunc)
		// progress.Watch()
		c.LinkTo(progress)
		c.Passthrough(event.RENDER, c.Next())

		for _, e := range c.data.All() {
			switch e.FsType() {
			case data.DirDatatype:

				if size, ok := c.cache.Get(e.FullPath()); ok {
					e.SetSize(size)
				} else {
					wg.Add(1)

					go func(e *data.FsEntry) {
						defer wg.Done()

						ctx.Observe(func() {
							size := c.scanSize(e.FullPath(), ctx)

							c.Publish("progress_message", pubsub.Message("Scanning "+e.FullPath()))

							e.SetSize(size)
							c.cache.Set(e.FullPath(), e.Size())
						})
					}(e)
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()
	c.Next().Unlink()
}

func (c *FsController) scanSize(path string, ctx *global.CancelContext) (total int) {
	if global.IsVirtualFs(path) {
		return
	}
	pathInfo, err := os.Stat(path)
	if err != nil {
		c.Publish("err", pubsub.Message(err.Error()))

		//	return
	}

	total += int(pathInfo.Size())

	contents, err := global.ReadIgnorePermission(path)
	if err != nil {
		c.Publish("err", pubsub.Message(err.Error()))

		//	return
	}

	for _, dirEntry := range contents {

		info, err := dirEntry.Info()
		if err != nil {
			c.Publish("progress_message", pubsub.Message("Warning "+err.Error()))

			//	c.printScanInfo("Warning: " + err.Error())

			continue
		} else {
			// c.Publish("progress_message", message.Message("Scanning "+info.Name()))
			if dirEntry.IsDir() {
				//	c.Publish("progress_message", message.Message("Scanning "+info.Name()))

				ctx.Observe(func() {
					size := c.scanSize(filepath.Join(path, info.Name()), ctx)

					//		c.Publish("progress_message", message.Message("Scanning "+filepath.Join(path, info.Name())))
					total += size
				})
			}

			total += int(info.Size())
		}

		// if err != nil {
		// 	// error = err
		// 	c.Broker.Publish("err", message.Message(err.Error()))
		// 	return
		// }

	}

	return
}

func (c *FsController) insertEntry(path string, info fs.FileInfo) {
	entry := c.pool.Get()

	switch true {

	// TODO: check for hard links also

	case global.IsSymlink(info):
		entry.SetFsType(data.SymlinkDatatype)

	case info.IsDir():
		//	entry.SetSize(int(info.Size()))
		entry.SetFsType(data.DirDatatype)

	default:
		entry.SetSize(int(info.Size()))
		entry.SetFsType(data.FileDatatype)

	}

	entry.SetName(filepath.Join(path, info.Name()))

	c.data.Insert(entry)
}
