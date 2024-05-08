package filesystem

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"

	"github.com/jbystronski/godirscan/pkg/app/data"
	"github.com/jbystronski/godirscan/pkg/global"
	g "github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

var (
	hideCursor         = g.HideCurson
	showCursor         = g.ShowCursor
	strlen             = g.StrLen
	cls                = g.ClearScreen
	cell               = g.Cell
	clear              = g.Clear
	rows               = g.Rows
	cols               = g.Cols
	trimEnd            = g.TrimEnd
	buildString        = g.BuildString
	fmtBold            = g.FmtBold
	themeMain          = g.ThemeMain
	ThemeAccent        = g.ThemeAccent
	themeBgHighlight   = g.ThemeBgHighlight
	themeHighlight     = g.ThemeHighlight
	themeBgHeader      = g.ThemeBgHeader
	themeHeader        = g.ThemeHeader
	themeBgSelect      = g.ThemeBgSelect
	themeBgPrompt      = g.ThemeBgPrompt
	themePrompt        = g.ThemePrompt
	themeSelect        = g.ThemeSelect
	truncate           = g.Truncate
	fmtPrompt          = g.FmtPrompt
	openCommandLine    = g.OpenCommandLine
	closeCommandLine   = g.CloseCommandLine
	receiveCommandLine = g.ReceiveCommandLine
	fmtSize            = g.FormatSize
)

const Space = g.Space

func (c *FsController) activeEntry() (*data.FsEntry, bool) {
	return c.data.Find(c.Index())
}

func (c *FsController) find(i int) (*data.FsEntry, bool) {
	return c.data.Find(i)
}

func (c *FsController) width() int {
	return c.panel.Width
}

func (c *FsController) contentLines() int {
	return c.panel.ContentLines()
}

func (c *FsController) outputFirstLine() int {
	return c.panel.OutputFirstLine()
}

func (c *FsController) outputLastLine() int {
	return c.panel.OutputLastLine()
}

func (c *FsController) contentWidth() int {
	return c.panel.ContentWidth()
}

func (c *FsController) contentLineStart() int {
	return c.panel.ContentStart()
}

// func rows() int {
// 	return termui.New().Rows()
// }

// func cols() int {
// 	return termui.New().Cols()
// }

func (c *FsController) getInput(prompt, placeholder string) string {
	clear(rows(), 1, cols())
	cmd := termui.NewCommandLine(rows(), 1, prompt, fmtPrompt, placeholder)

	input := cmd.WaitInput()

	// printHelper()

	return input
}

func tryCreateSymlink(srcPath, targetPath string) (bool, error) {
	fileInfo, err := os.Lstat(srcPath)
	if err != nil {
		return false, err
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		symlinkTarget, linkErr := os.Readlink(srcPath)

		if linkErr != nil {
			return false, linkErr
		}

		symlinkErr := os.Symlink(symlinkTarget, targetPath)
		if symlinkErr != nil {
			return false, symlinkErr
		}

		return true, nil

	}
	return false, nil
}

func (c *FsController) UpdateSize() {
	c.panel.Top = 2
	c.panel.Height = rows() - 2
	c.panel.Width = cols() / 2
	c.panel.SetPadding(0, 0, 1, 0)

	// p.Main.Width = cols() / 2
	// p.HeaderBox.Width = cols() / 2

	// p.HeaderBox.OffsetBottom = rows() - 1

	if c.panel.Left > 1 {
		c.panel.Left = cols()/2 + 1
	} else {
		c.panel.Left = 1
	}

	c.MinOffset = c.panel.OutputFirstLine()
	c.MaxOffset = c.panel.OutputLastLine() - 1

	// if p.Main.OffsetLeftStart > 1 {
	// 	p.Main.OffsetLeftStart = p.Main.Width + 1
	// 	p.HeaderBox.OffsetLeftStart = p.Main.Width + 1
	// }
}

func (c *FsController) updateTotalSize() {
	cell(rows()-2, c.panel.ContentStart())
	fmt.Print(totalSize(c.data.Size()))
}

func (c *FsController) sendError(err error) {
	c.Broker.Publish("err", pubsub.Message(err.Error()))
}

func (c *FsController) updateParentStoreSize(updatedSize int) {
	path := c.root

	for path != global.GetRootDirectory() {
		if parentDir, ok := global.GetParentDirectory(path); ok {
			if size, ok := c.cache.Get(parentDir); ok {
				c.cache.Set(parentDir, size+updatedSize)
				path = parentDir

				continue
			}
			return
		}
	}
}

func (c *FsController) Create(name string, en *data.FsEntry) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("name can't be empty")
	}

	if strings.ContainsAny(name, string(os.PathSeparator)) {
		return fmt.Errorf("%s \"%v\"", "Name cannot contain path separator", string(os.PathSeparator))
	}

	path := filepath.Join(c.root, name)

	switch en.FsType() {
	case data.DirDatatype:
		err := os.Mkdir(path, 0o777)
		if err != nil {
			return err
		}

	case data.FileDatatype:
		_, err := os.Create(path)
		if err != nil {
			return err
		}

	}

	return nil
}

func (c *FsController) ObserveTicker(ctxDone <-chan struct{}, interval time.Duration, intervalCallback func()) (*sync.WaitGroup, chan struct{}) {
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
			//	common.Log("ticker has been stopped")
			//	wg.Done()
		}()
		for {
			select {

			// case err := <-c.internalErrChan:

			// 	c.ErrorChan <- err
			// 	return

			case <-ticker.C:
				if init {
					//	common.Log("TICK")
					intervalCallback()
				}

			case <-ctxDone:
				return

			case <-done:
				//	common.Log("stoppng ticker ib observer")
				return
			}
		}
	}()

	return &wg, done
}

func clearPrompt() {
	clear(rows(), 1, cols())
}

func printInfo(s string) {
	cell(rows(), 1)
	fmt.Print(s)
}

func ExecuteFile(path string, executors map[string]string) error {
	ext := filepath.Ext(path)

	var cmd *exec.Cmd

	if customCmd, ok := executors[ext]; ok {
		cmd = exec.Command(customCmd, path)
	} else {
		switch runtime.GOOS {
		case "darwin":
			{
				cmd = exec.Command("open", path)
			}
		case "windows":
			{
				cmd = exec.Command("cmd", "/c", "start", path)
			}
		default:
			{
				cmd = exec.Command("xdg-open", path)
			}
		}
	}

	cls()
	fmt.Println(ThemeAccent() + "Press esc to return, command execution output: " + "\033[0m")
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	for {

		_, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		if key == keyboard.KeyEsc {
			return nil
		}

	}
}

func header(h string) string {
	return fmtBold(themeBgHeader(), themeHeader(), h)
}

func file(sep, line string) string {
	return fmtBold(sep, ThemeAccent(), line)
}

func activeRow(sep, line string) string {
	return fmtBold(themeBgHighlight(), themeHighlight(), sep, line)
}

func directory(sep, line string) string {
	return fmtBold(sep, themeMain(), line)
}

func searchResult(sep string, en data.FsEntry) string {
	return fmtBold(sep, ThemeAccent(), en.FullPath())
}

func selectedRow(sep, line string) string {
	return fmtBold(themeBgSelect(), themeSelect(), sep, line)
}

func totalSize(s int) string {
	return fmtBold(themeBgHeader(), themeHeader(), printSizeAsString(s))
}

func printSizeAsString(size int) string {
	return fmt.Sprintf("%v", fmtSize(size))
}

func symlink(sep, line string) string {
	return directory(sep, line)
}
