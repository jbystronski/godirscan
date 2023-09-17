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

	fs "github.com/jbystronski/godirscan/pkg/filesystem"
)

var actionMap = map[k.Key]c.ControllerAction{
	k.KeyEnter:      c.Execute,
	k.KeyCtrlE:      c.ExecuteCmd,
	k.KeyDelete:     c.Delete,
	k.KeyF8:         c.Delete,
	k.KeyCtrlR:      c.Rename,
	k.KeyF7:         c.CreateDirectory,
	k.KeyCtrlW:      c.CreateFile,
	k.KeyCtrlF:      c.Find,
	k.KeyCtrlS:      c.Sort,
	k.KeyF6:         c.Move,
	k.KeyCtrlV:      c.Copy,
	k.KeyF9:         c.Scan,
	k.KeyHome:       c.MoveToTop,
	k.KeyEnd:        c.MoveToBottom,
	k.KeyPgdn:       c.PageDown,
	k.KeyPgup:       c.PageUp,
	k.KeyArrowLeft:  c.MoveLeft,
	k.KeyArrowRight: c.MoveRight,
	k.KeyInsert:     c.Select,
	k.KeyCtrlA:      c.SelectAll,
	k.KeyCtrlG:      c.GoTo,
	k.KeyArrowDown:  c.MoveDown,
	k.KeyArrowUp:    c.MoveUp,
	k.KeyCtrl4:      c.Edit,
}

const (
	topMargin    = 2
	bottomMargin = 3
	borderWidth  = 1
)

var (
	ctrl, alt *fs.FsController

	resizeChan = make(chan os.Signal, 1)

	errorChan = make(chan error, 1)

	selected = fs.NewSelected()
)

func init() {
	os.Truncate("/home/kb/log", 0)
	c.Log("start")
	c.ParseConfigFile(c.Cfg)
	c.ParseColorSchema(c.Cfg.CurrentSchema, c.CurrentTheme)
}

func navigate() {
	defer func() {
		if r := recover(); r != nil {

			c.FlashError(fmt.Errorf("%s", r))
			time.Sleep(time.Second * 2)
			navigate()
		}
	}()

	keysEvents, err := k.GetKeys(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := k.Close()
		if err != nil {
			c.FlashError(err)
		}
	}()

	signal.Notify(resizeChan, syscall.SIGWINCH)

	for {
		select {

		case <-ctrl.Tick():
			if ctrl.Tickable.IsInitialized() {
				c.Log("tick")

				ctrl.Map(c.Render)
			}

		case <-alt.Tick():
			if alt.Tickable.IsInitialized() {
				c.Log("alt tick")

				alt.Map(c.Render)
			}

		case <-ctrl.Done:
			common.Log("Done calculating size")
			ctrl.Map(c.Refresh)

		case <-alt.Done:
			alt.Map(c.Refresh)

		case err := <-errorChan:

			panic(err)

		case <-resizeChan:
			c.ClearScreen()
			ctrl.Map(c.Resize)
			alt.Map(c.Resize)
			continue

		case event := <-keysEvents:

			switch event.Key {

			case k.KeyCtrl5:

				func() {
					c.ClearScreen()

					c.PrintBanner(*c.CurrentTheme)
					fmt.Println("total lines: ", ctrl.TotalLines())

					fmt.Println("output lines: ", ctrl.Lines())
					fmt.Println("output first line: ", ctrl.OutputFirstLine())
					fmt.Println("output last line: ", ctrl.OutputLastLine())
					fmt.Println("current index ", ctrl.Index())

					fmt.Println("first index: ", ctrl.ChunkStart())
					fmt.Println("last index:", ctrl.ChunkEnd())

					ctrl.FsNavigator.Print()

					for {
						_, key, err := k.GetKey()
						if err != nil {
							fmt.Println(err)
						}

						if key == k.KeyEsc {
							return
						}

					}
				}()

				// common.PrintHelp(*c.CurrentTheme, *ctrl)
				ctrl.Map(c.Refresh)
				alt.Map(c.Refresh)
				continue

			case k.KeyCtrlSlash:
				c.Cfg.ChangeTheme()

				ctrl.Map(c.Refresh)
				alt.Map(c.Refresh)

			case k.KeyCtrlH:
				common.ClearScreen()
				ctrl.SetWidth(c.PaneWidth() * 2)
				ctrl.SetOffsetLeftStart(1)
				ctrl.Map(c.Refresh)

			case k.KeyTab:

				ctrl.SetActive(false)
				alt.SetActive(true)

				ctrl, alt = alt, ctrl

				ctrl.Map(c.Render)
				alt.Map(c.Render)
				continue

			case k.KeyCtrlD:
				c.ViewLog()
				ctrl.Map(c.Refresh)
				alt.Map(c.Refresh)

			case k.KeyEsc, k.KeyF10, k.KeyCtrlC:
				c.ClearScreen()

				return

			default:

				if action, ok := actionMap[event.Key]; ok {
					ctrl.Map(action)
				} else {
					continue
				}

			}

		}
	}
}

func main() {
	defer func() {
		c.ClearScreen()
		fmt.Print("\033[?25h")
		if k.IsStarted(time.Millisecond * 50) {
			k.Close()
		}
	}()

	startDir := startDirectory()

	if startDir == "" {
		return
	}

	var err error

	vb := c.ViewBox{}
	vb.SetOffsetTop(topMargin)
	vb.SetOffsetBottom(bottomMargin)
	vb.SetTotalLines(c.NumVisibleLines())
	vb.SetWidth(c.PaneWidth())
	vb.SetOffsetTopStart(2)
	vb.SetOffsetLeftStart(1)
	vb.SetPadding(0)
	vb.SetTheme(c.CurrentTheme)

	ctrl, err = fs.NewController(errorChan, startDir, fs.ViewBox{
		ViewBox: vb,
	}, *selected)

	ctrl.SetActive(true)

	if err != nil {
		c.FlashError(err)
		return
	}

	vb2 := c.ViewBox{}
	vb2.SetOffsetTop(topMargin)
	vb2.SetOffsetBottom(bottomMargin)
	vb2.SetTotalLines(c.NumVisibleLines())
	vb2.SetWidth(c.PaneWidth())
	vb2.SetOffsetTopStart(2)
	vb2.SetOffsetLeftStart(c.PaneWidth() + 1)
	vb2.SetPadding(0)
	vb2.SetTheme(c.CurrentTheme)

	alt, err = fs.NewController(errorChan, startDir, fs.ViewBox{
		ViewBox: vb2,
	}, *selected)

	if err != nil {
		c.FlashError(err)
		return
	}

	ctrl.Alt = alt
	alt.Alt = ctrl

	navigate()
}

func startDirectory() string {
	var path string
	_ = k.Open()

	defer func() {
		k.Close()
	}()

	c.ClearScreen()

	c.PrintBanner(*c.CurrentTheme)

	path = c.WaitInput("Scan directory: ", c.Cfg.DefaultRootDirectory, c.Coords{Y: 10, X: 1}, errorChan)

	c.ClearScreen()

	return path
}
