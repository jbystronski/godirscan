package controller

import (
	"context"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem"
	fs "github.com/jbystronski/godirscan/pkg/filesystem"
	v "github.com/jbystronski/godirscan/pkg/filesystem/viewbox"
	"github.com/jbystronski/godirscan/pkg/navigator"
)

type Controller struct {
	config *common.Config

	ctx *common.CancelContext

	resizeContext context.Context

	active   bool
	cache    *fs.CacheStore
	pool     *fs.FsDataPool
	selected *fs.Selected
	path     string

	navigator.ChunkNavigator
	fs.DataAccessor

	v.FsViewBox

	ErrorChan       chan<- (error)
	internalErrChan <-chan (error)

	Alt *Controller

	backtrace   common.MapAccessor[string, int]
	defaultSort int64
	hidden      bool
	actionMap   map[k.Key]func()
}

const (
	sortByName = 1
	osrtBySizeAsc
	sortBySizeDesc
	sortByType
)

func NewController(cache *filesystem.CacheStore, selected *filesystem.Selected, dataPool *filesystem.FsDataPool, errChan chan<- error, vb v.FsViewBox) *Controller {
	c := &Controller{
		cache:    cache,
		selected: selected,
		pool:     dataPool,
		ctx:      &common.CancelContext{},

		ErrorChan:       errChan,
		internalErrChan: make(chan error),
		DataAccessor:    fs.NewFsData(),

		ChunkNavigator: *navigator.NewChunkNavigator(),
		FsViewBox:      vb,
		backtrace:      &common.GenericMap[string, int]{},
		hidden:         false,
		defaultSort:    1,
	}

	c.actionMap = map[k.Key]func(){
		// k.KeyEnter:      c.execute,
		// k.KeyCtrlE:      c.executeCmd,
		// k.KeyDelete:     c.delete,
		// k.KeyF8:         c.delete,
		// k.KeyCtrlR:      c.rename,
		// k.KeyF7:         c.newDirectory,
		// k.KeyCtrlW:      c.newFile,
		// k.KeyCtrlF:      c.search,
		// k.KeyCtrlS:      c.sort,
		// k.KeyF6:         c.Move,
		// k.KeyCtrlV:      c.Copy,
		// k.KeyF9:         c.scanDir,
		// k.KeyHome:       c.top,
		// k.KeyEnd:        c.bottom,
		// k.KeyPgdn:       c.pgDown,
		// k.KeyPgup:       c.pgUp,
		// k.KeyArrowLeft:  c.left,
		// k.KeyArrowRight: c.right,
		// k.KeyInsert:     c.selectEntry,
		// k.KeyCtrlA:      c.selectAllEntries,
		// k.KeyCtrlG:      c.goTo,
		// k.KeyArrowDown:  c.down,
		// k.KeyArrowUp:    c.up,
		//	k.KeyCtrl4:      c.edit,
		//	k.KeyTab:        c.changeController,
		k.KeyCtrlH: c.hideAlt,
		//	k.KeyEsc:   c.ctx.Cancel,
		//	k.KeyCtrlSlash:  c.updateTheme,
		k.KeyCtrl5: c.printHelp,
		//	k.KeyF5:         c.Refresh,
	}

	c.SetChunkLines(c.Lines())

	c.PrintBox()
	return c
}

func (c *Controller) Active() bool {
	return c.active
}

func (c *Controller) SetActive(active bool) {
	c.active = active
}

func (c *Controller) SetPath(path string) {
	c.path = path
}
