package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
	c "github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem/controller"
)

func init() {
	c.ParseConfigFile(c.Cfg)
	c.ParseColorSchema(c.Cfg.CurrentSchema, c.CurrentTheme)
}

func main() {
	defer func() {
		common.ClearScreen()
		fmt.Print("\033[?25h")
		err := k.Close()
		if err != nil {
			c.FlashError(err)
		}
	}()

	var (
		resizeChan  = make(chan os.Signal, 1)
		mainErrChan = make(chan error, 1)
		resumeChan  = make(chan struct{}, 1)
		initChan    = make(chan struct{}, 1)
		ctrl        *controller.Controller
	)

	initChan <- struct{}{}

	keysEvents, err := k.GetKeys(1)
	if err != nil {
		panic(err)
	}

	signal.Notify(resizeChan, syscall.SIGWINCH)

	for {
		select {

		case <-initChan:
			common.ClearScreen()
			ctrl = controller.Init(mainErrChan, c.Cfg)

		case err := <-mainErrChan:

			//	panic(err)
			c.FlashError(fmt.Errorf("%s%s", "error occured", err))
			time.Sleep(time.Second * 2)
			resumeChan <- struct{}{}

		case <-resizeChan:
			common.ClearScreen()
			ctrl.Refresh()
			ctrl.Alt.Refresh()

		case <-resumeChan:
			continue

		case event := <-keysEvents:

			switch event.Key {

			case k.KeyCtrlC, k.KeyCtrlQ, k.KeyF10:

				return

			default:
				ctrl.MapKey(event.Key)

			}

		}
	}
}
