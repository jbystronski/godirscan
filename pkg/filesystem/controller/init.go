package controller

import (
	"time"

	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem"
	"github.com/jbystronski/godirscan/pkg/viewbox"

	fsViebox "github.com/jbystronski/godirscan/pkg/filesystem/viewbox"
)

func Init(mainErrChan chan error, config *common.Config) *Controller {
	selected := filesystem.NewSelected()
	cache := filesystem.NewCacheStore(31457280)
	dataPool := filesystem.NewDataPool()
	var ctrl, alt *Controller

	const (
		topMargin    = 2
		bottomMargin = 3
		borderWidth  = 1
	)

	vb := viewbox.ViewBox{}
	vb.SetOffsetTop(topMargin)
	vb.SetOffsetBottom(bottomMargin)
	vb.SetTotalLines(common.NumVisibleLines())
	vb.SetWidth(common.PaneWidth())
	vb.SetOffsetTopStart(2)
	vb.SetOffsetLeftStart(1)
	vb.SetPadding(0)
	vb.SetTheme(common.CurrentTheme)

	vb.PrintBanner()
	startDirectory := common.WaitInput("Scan directory: ", config.DefaultRootDirectory, common.Coords{Y: 11, X: 1}, mainErrChan)

	common.ClearScreen()

	ctrl = NewController(cache, selected, dataPool, mainErrChan, fsViebox.FsViewBox{
		ViewBox: vb,
	})

	ctrl.SetActive(true)

	vb2 := vb

	vb2.SetOffsetLeftStart(common.PaneWidth() + 1)

	alt = NewController(cache, selected, dataPool, mainErrChan, fsViebox.FsViewBox{
		ViewBox: vb2,
	})

	ctrl.config = config
	alt.config = config

	ctrl.Alt, alt.Alt = alt, ctrl

	ctrl.SetStore(startDirectory)

	// go func() {
	time.Sleep(time.Millisecond * 600)
	alt.SetStore(startDirectory)
	// }()

	return ctrl
}
