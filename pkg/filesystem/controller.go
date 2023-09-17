package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
)

type FsController struct {
	common.Controller
	common.Cancelable
	common.Tickable
	FsNavigator
	Scanner
	Finder
	FsStoreAccessor
	ViewBox
	selected  *Selected
	ErrorChan chan<- (error)

	Done chan struct{}
	Alt  *FsController

	backtraceMap map[string]int
	defaultSort  int
}

const (
	sortByName = 1
	osrtBySizeAsc
	sortBySizeDesc
	sortByType
)

func NewController(errChan chan<- error, path string, vb ViewBox, selected Selected) (*FsController, error) {
	dirReader := &DirReader{}

	cancelCtx := &common.CancelCtx{}
	ticker := &common.Ticker{}
	ticker.SetInterval(time.Second * 1)

	c := &FsController{
		Done:       make(chan struct{}),
		Controller: &common.BaseController{},
		Tickable:   ticker,
		Cancelable: cancelCtx,
		ErrorChan:  errChan,
		Scanner:    Scanner{DirReader: *dirReader, Cancelable: cancelCtx, Tickable: ticker},

		Finder: Finder{DirReader: *dirReader, Cancelable: cancelCtx, Tickable: ticker},
		FsStoreAccessor: &FsStore{
			StoreAccessor: &common.Store{
				DataAccessor: &Entries{},
			},
		},
		selected:    &selected,
		FsNavigator: FsNavigator{ChunkNavigator: common.ChunkNavigator{BaseNavigator: common.BaseNavigator{}}},
		ViewBox:     vb,

		defaultSort: 1,
	}

	vb.Trimmer = common.Trimmer{}
	c.SetChunkLines(c.Lines())
	err := c.SetStore(path)
	if err != nil {
		return nil, err
	}

	c.SetChunk(c.Data().Len())

	c.SetActionMap(map[common.ControllerAction]func(){
		common.MoveDown:        c.nextEntry,
		common.MoveUp:          c.prevEntry,
		common.MoveLeft:        c.goUpDirectory,
		common.MoveRight:       c.enterDirectory,
		common.MoveToTop:       c.firstEntry,
		common.MoveToBottom:    c.lastEntry,
		common.PageDown:        c.moveDownEntries,
		common.PageUp:          c.moveUpEntries,
		common.Delete:          c.deleteSelected,
		common.Scan:            c.scanDirectory,
		common.Execute:         c.executeEntry,
		common.ExecuteCmd:      c.executeCmd,
		common.Rename:          c.renameEntry,
		common.CreateDirectory: c.newDirectory,
		common.CreateFile:      c.newFile,
		common.Sort:            c.sortEntries,
		common.Edit:            c.editEntry,
		common.Select:          c.selectEntry,
		common.SelectAll:       c.selectAllEntries,
		common.GoTo:            c.goToEntry,
		common.Resize:          c.resizeTerminal,
		common.Copy:            c.copyEntries,
		common.Move:            c.moveEntries,
		common.Find:            c.findEntries,

		common.Render:  c.render,
		common.Refresh: c.fullRender,
	})
	c.Data().sortByType()
	c.fullRender()

	return c, nil
}

func (c *FsController) store() FsStoreAccessor {
	return c.FsStoreAccessor
}

func (c *FsController) current() *FsFiletype {
	return c.Data().Find(c.Index())
}

func (c *FsController) get(index int) *FsFiletype {
	return c.Data().Find(index)
}

func (c *FsController) len() int {
	return c.Data().Len()
}

func (c *FsController) SetStore(dir string) error {
	var error error

	c.ResolveUserDirectory(&dir)

	c.store().SetName(dir)

	entries, err := c.scan(dir)
	if err != nil {

		error = err
		return error
	}

	c.SetData(entries)

	c.SetChunk(c.Data().Len())

	// go func() {
	// 	err := c.scanDataSize(c.Data().All())
	// 	if err != nil {
	// 		error = err
	// 		return
	// 	}
	// 	common.Log("done calucalting of ", c.store().Name(), " store size is ", c.store().Size())
	// 	c.Done <- struct{}{}
	// }()

	return error
}

func (c *FsController) DefaultSort() {
	switch c.defaultSort {
	case 1:
		c.Data().SortByName()
	case 2:
		sort.Sort(c.Data())
	case 3:
		sort.Sort(sort.Reverse(c.Data()))
	case 4:
		c.Data().sortByType()
	default:
		c.Data().SortByName()
	}
}

func (c *FsController) SetDefaultSort(s int) {
	c.defaultSort = s
}

func (c *FsController) moveUpEntries() {
	c.MovePgUp(c.len())

	c.render()
}

func (c *FsController) moveDownEntries() {
	c.MovePgDown(c.len())

	c.render()
}

func (c *FsController) nextEntry() {
	c.MoveDown(c.len())

	c.render()
}

func (c *FsController) prevEntry() {
	c.MoveUp(c.len())

	c.render()
}

func (c *FsController) wrapInput(prompt, defaultOpt string) string {
	answ := common.WaitInput(prompt, defaultOpt, c.GoToPromptCell(), c.ErrorChan)

	defer func() {
		c.PrintBox()
		c.Alt.PrintBox()
		c.printTotalSize(c.store().Size())
		c.printTotalSize(c.Alt.store().Size())
	}()

	return answ
}

func (c *FsController) scanDirectory() {
	dir := c.wrapInput("Scan directory: ", common.Cfg.DefaultRootDirectory)

	if dir == "" {
		return
	}

	c.Reset()

	err := c.SetStore(dir)
	if err != nil {
		c.ErrorChan <- err
	}
	c.DefaultSort()
	c.fullRender()
}

func (c *FsController) goUpDirectory() {
	dir, ok := c.store().ParentStoreName()

	if !ok {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		c.ErrorChan <- err
	}

	if c.Backtrace(dir, len(entries)) {
		err := c.SetStore(dir)
		if err != nil {
			c.ErrorChan <- err
		}

		c.DefaultSort()
		c.fullRender()

	}
}

func (c *FsController) enterDirectory() {
	if c.len() == 0 {
		return
	}

	switch (*c.current()).(type) {
	case *FsDirectory:

		c.SetBacktrace(c.store().Name())
		err := c.SetStore((*c.current()).FullPath())
		if err != nil {
			c.ErrorChan <- err
		}
		c.DefaultSort()
		c.Reset()
		c.SetChunk(c.Data().Len())
		c.fullRender()
	}
}

func (c *FsController) firstEntry() {
	c.MoveToTop(c.len())

	c.render()
}

func (c *FsController) lastEntry() {
	c.MoveToBottom(c.len())

	c.render()
}

func (c *FsController) deleteSelected() {
	answ := c.wrapInput("Delete selected entries", "y")

	if answ == "y" || answ == "Y" {

		doneChan := make(chan struct{})
		messageChan := make(chan string)

		go func() {
			common.PrintProgress(doneChan, messageChan, 1, common.NumVisibleLines())
		}()

		ok, err := c.selected.Delete(c.store().Name(), messageChan)
		if err != nil {
			doneChan <- struct{}{}
			c.ErrorChan <- err
		}

		if ok {

			err := c.SetStore(c.store().Name())
			if err != nil {
				c.ErrorChan <- err
			}

			err = c.Alt.SetStore(c.Alt.store().Name())
			if err != nil {
				c.ErrorChan <- err
			}

			if c.len() == 0 {
				c.Reset()
			} else if c.Index() > c.len()-1 {
				c.SetIndex(c.len() - 1)
			}

			c.fullRender()
			c.Alt.fullRender()
		} else {
			c.render()
			c.Alt.render()
		}

		doneChan <- struct{}{}

	}
}

func (c *FsController) findEntries() {
	answ := c.wrapInput("Find by (1) pattern (2) size", "")

	if answ != "1" && answ != "2" {
		return
	}
	path := c.wrapInput("Path to search from", c.store().Name())

	if path == "" {
		return
	}

	// if info, err := c.GetPathInfo(path); err != nil {

	// 	c.ErrorChan <- err

	// 	if !info.IsDir() {
	// 		c.ErrorChan <- errors.New("path must be a directory")
	// 	}

	// 	return

	// }

	pattern := c.wrapInput("Search pattern", "")

	c.SetData(&Entries{})

	if answ == "1" {

		test := func(info fs.FileInfo) bool {
			return true
		}

		c.find(path, pattern, c.ErrorChan, c.Done, c.Data().Insert, test)
	}

	if answ == "2" {
		answ := c.wrapInput("Select unit (0) bytes (1) kb (2) mb (3) gb", "")

		if answ != "0" && answ != "1" && answ != "2" && answ != "3" {
			return
		}

		unit, _ := strconv.Atoi(answ)

		unitName := common.StorageUnits[unit]

		answ = c.wrapInput(fmt.Sprintf("%s %s", "Type min value in", unitName), "0")

		if answ == "" {
			return
		}

		min, err := strconv.ParseFloat(answ, 64)
		if err != nil {
			c.ErrorChan <- err
			return
		}

		if min < 0 {
			min = 0
		}

		answ = c.wrapInput(fmt.Sprintf("%s %s", "Type max value ( no value means unlimited ) in ", unitName), "")

		if answ == "" {
			answ = "0"
		}

		max, err := strconv.ParseFloat(answ, 64)
		if err != nil {
			c.ErrorChan <- err
			return
		}

		if max < 0 || min < 0 {
			c.ErrorChan <- errors.New("negative values are not allowed")
			return
		}

		if max > 0 && max < min {
			c.ErrorChan <- fmt.Errorf("max value: %v, can't be lower than min value: %v", max, min)
			return
		}

		minV, maxV := common.ToBytes(unit, min, max)

		test := func(info fs.FileInfo) bool {
			if info.Size() >= minV {
				if maxV != 0 && info.Size() <= maxV || maxV == 0 {
					return true
				}
			}

			return false
		}

		c.find(path, pattern, c.ErrorChan, c.Done, c.Data().Insert, test)

	}
}

func (c *FsController) executeEntry() {
	entry := *c.current()

	if executable, ok := entry.(executable); ok {

		err := executable.execute()
		if err != nil {
			c.ErrorChan <- err
		}

		c.fullRender()
		c.Alt.fullRender()

	} else if _, ok := entry.(*FsDirectory); ok {
		c.enterDirectory()
	}
}

func (c *FsController) executeCmd() {
	input := c.wrapInput("Run command: ", "")

	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	//	cmd.Stdin = os.Stdin
	common.ClearScreen()
	fmt.Println("Press esc to return, command execution output: " + "\033[0m")
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		c.ErrorChan <- err
	}

	for {

		_, key, err := keyboard.GetKey()
		if err != nil {
			c.ErrorChan <- err
		}

		if key == keyboard.KeyEsc {
			break
		}

	}

	common.ClearScreen()
	c.fullRender()
	c.Alt.fullRender()
}

func (c *FsController) renameEntry() {
	if c.len() > 0 {

		en := c.current()
		oldPath := (*en).FullPath()
		newName := c.wrapInput("Rename:", (*en).Name())

		ok, err := (*en).Rename(newName)
		c.ErrorChan <- err

		if ok {

			c.render()
			if c.store().Name() == c.Alt.store().Name() {
				match := c.Data().FindByPath(oldPath)
				(*match).SetName(newName)
				c.Alt.fullRender()
			}
		}

	}
}

func (c *FsController) newDirectory() {
	newName := c.wrapInput("Create directory: ", "")

	ok := c.Create(newName, &FsDirectory{})

	if ok {

		err := c.SetStore(c.store().Name())
		if err != nil {
			c.ErrorChan <- err
		}

		c.fullRender()
		err = c.Alt.SetStore(c.Alt.store().Name())
		if err != nil {
			c.ErrorChan <- err
		}

		c.Alt.fullRender()
	} else {
		c.render()
		c.Alt.render()
	}
}

func (c *FsController) newFile() {
	newFileName := c.wrapInput("Create a new file: ", "")

	path := filepath.Join(c.Name(), newFileName)

	_, err := os.Stat(path)

	if err == nil {
		answ := c.wrapInput(fmt.Sprintf("%s (%s) %s", "File", newFileName, "already exists, do you wish to override it?"), "n")
		if answ == "y" || answ == "Y" {
			c.truncate(path)
		} else {
			c.render()
			c.Alt.render()

		}

	} else {

		ok := c.Create(newFileName, &FsFile{})

		if ok {

			err := c.SetStore(c.store().Name())
			if err != nil {
				c.ErrorChan <- err
			}

			err = c.Alt.SetStore(c.Alt.store().Name())

			if err != nil {
				c.ErrorChan <- err
			}

			c.fullRender()
			c.Alt.fullRender()
		} else {
			c.render()
			c.Alt.render()
		}
	}
}

func (c *FsController) sortEntries() {
	answ := c.wrapInput("Sort by (1) Name (2) Size ASC (3) Size DESC (4) Type", "")

	if answ != "1" && answ != "2" && answ != "3" && answ != "4" {
		return
	}

	s, _ := strconv.Atoi(answ)

	c.SetDefaultSort(s)

	switch answ {
	case "1":
		c.Data().SortByName()

	case "2":
		sort.Sort(c.Data())

	case "3":
		sort.Sort(sort.Reverse(c.Data()))

	case "4":
		c.Data().sortByType()

	default:
		c.Data().SortByName()

	}

	c.fullRender()
}

func (c *FsController) editEntry() {
	if c.len() > 0 {
		if editable, ok := (*c.current()).(Editable); ok {

			common.ClearScreen()
			err := editable.Edit()
			if err != nil {
				c.ErrorChan <- err
			}

			common.HideCursor()
			c.fullRender()
			c.Alt.fullRender()
		}
	}
}

func (c *FsController) selectEntry() {
	if c.len() > 0 {
		c.selected.DumpPrevious(c.store().Name())

		en := c.current()

		//	en := (*e).(FsFiletype)

		c.selected.Select((*en).FullPath())
		if c.Alt.store().Name() == c.store().Name() {
			c.Alt.fullRender()
		}
		c.nextEntry()

	}
}

func (c *FsController) selectAllEntries() {
	if c.len() > 0 {
		c.selected.DumpPrevious(c.store().Name())
		c.selected.SelectAll(c.Data().All())
		c.fullRender()
		if c.Alt.store().Name() == c.store().Name() {
			c.Alt.fullRender()
		}
	}
}

func (c *FsController) goToEntry() {
	if c.len() > 0 {
		var matches []int
		var curr int

		test := c.wrapInput("Go to: ", "")
		c.render()
		c.Alt.render()

		if test == "" {
			return
		}

		for index, en := range c.Data().All() {
			if strings.Contains((*en).Name(), test) {
				matches = append(matches, index)
			}
		}

		if len(matches) > 0 {

			c.JumpTo(c.PrevIndex(), matches[0], c.len())

			c.render()
			c.GoToCell(c.TotalLines()-2, c.ContentLineStart())

		} else {
			c.ErrorChan <- errors.New("no matches found")

			return
		}

	MatchLoop:
		for {

			_, key, err := keyboard.GetKey()
			if err != nil {
				c.ErrorChan <- err
				break
			}
			switch key {

			case keyboard.KeyInsert:

				en := c.current()

				file := (*en).(FsFiletype)

				c.selected.DumpPrevious(c.store().Name())
				c.selected.Select((file).FullPath())
				c.render()

			case keyboard.KeyArrowUp, keyboard.KeyArrowLeft:
				if curr == 0 {
					continue
				}
				prev := matches[curr]
				curr--
				c.JumpTo(prev, matches[curr], c.len())

				c.render()
				common.Cell(c.TotalLines()-2, c.ContentLineStart())

			case keyboard.KeyArrowDown, keyboard.KeyArrowRight:
				if curr == len(matches)-1 {
					continue
				}

				prev := matches[curr]
				curr++
				c.JumpTo(prev, matches[curr], c.len())

				c.render()
				common.Cell(c.TotalLines()-2, c.ContentLineStart())

			case keyboard.KeyEsc, keyboard.KeyEnter:
				common.ClearLine()
				common.CarriageReturn()
				break MatchLoop
			}

		}

		c.render()
		c.Alt.render()

	}
}

func (c *FsController) resizeTerminal() {
	c.SetTotalLines(common.NumVisibleLines())
	c.SetChunkLines(c.Lines())
	c.SetChunk(c.len())

	c.SetWidth(common.PaneWidth())

	if c.OffsetLeftStart() > 1 {
		c.SetOffsetLeftStart(common.PaneWidth() + 1)
	}

	c.fullRender()
}

func (c *FsController) moveEntries() {
	if c.selected.IsEmpty() {
		return
	}

	answ := c.wrapInput(fmt.Sprintf("%s", "Move selected into the current directory? :"), "y")
	c.render()

	if answ != "y" {
		return
	}

	relocate(common.Move, c.ErrorChan, c.selected.All(), c.store().Name())
}

func (c *FsController) copyEntries() {
	if c.selected.IsEmpty() {
		return
	}

	answ := c.wrapInput(fmt.Sprintf("%s", "Copy selected into the current directory? :"), "y")
	c.render()

	if answ != "y" {
		return
	}

	relocate(common.Copy, c.ErrorChan, c.selected.All(), c.store().Name())
}

func (c *FsController) Create(name string, ftype FsFiletype) bool {
	var err error
	name = strings.TrimSpace(name)

	if name == "" {
		return false
	}

	if strings.ContainsAny(name, string(os.PathSeparator)) {
		c.ErrorChan <- fmt.Errorf("%s \"%v\"", "Name cannot contain", string(os.PathSeparator))
		return false
	}

	path := filepath.Join(c.store().Name(), name)

	if _, ok := ftype.(*FsDirectory); ok {
		err = os.Mkdir(path, 0o777)

		return true
	}

	if _, ok := ftype.(*FsFile); ok {
		_, err = os.Create(path)
	}
	if err != nil {
		c.ErrorChan <- err
		return false
	}

	return true
}

func (c *FsController) truncate(path string) {
	err := os.Truncate(path, 0)
	if err != nil {
		c.ErrorChan <- err
	}
}

// func (c *FsController) vfullRender() {
// 	c.PrintBox()

// 	c.ClearLine(1, c.ContentLineStart()-1, c.ContentWidth())

// 	c.printHeader(c.store().Name())

// 	if c.len() == 0 {
// 		c.printEmptyFolder(c.Active())
// 		return
// 	}

// 	c.GoToCell(c.Line(c.Index()), c.ContentLineStart())

// 	line := c.FirstLine()

// 	for i := c.ChunkStart(); i <= c.ChunkEnd(); i++ {
// 		file := c.get(i)

// 		if file == nil {
// 			continue
// 		}

// 		opts := UpdateLineOptions{
// 			en:                 *file,
// 			index:              i,
// 			isActive:           file == c.current(),
// 			isControllerActive: c.Active(),
// 			isSelected:         c.selected.IsSelected((*file).FullPath()),
// 			isLast:             i == c.Last(),
// 		}

// 		c.UpdateLine(opts)

// 		line++
// 		c.GoToCell(line, c.ContentLineStart())
// 	}

// 	for line <= c.OutputLastLine() {
// 		c.ClearLine(line, c.ContentLineStart(), c.ContentWidth())
// 		line++
// 	}

// 	c.GoToCell(c.TotalLines()-2, c.ContentLineStart())
// 	fmt.Print(printSizeAsString(c.store().Size()))

// 	c.GoToCell(c.TotalLines(), 1)
// }

func (c *FsController) Selected() int {
	return len(c.selected.All())
}

func (c *FsController) fullRender() {
	opts := RenderOpts{
		start:         c.ChunkStart(),
		end:           c.ChunkEnd(),
		storeSize:     c.store().Size(),
		index:         c.Index(),
		entries:       c.Data().All(),
		checkSelected: c.selected.IsSelected,
		selectKey: func(f FsFiletype) string {
			return f.FullPath()
		},
		storeName:          c.store().Name(),
		isControllerActive: c.Active(),
	}

	c.ViewBox.RenderAll(opts)
}

func (c *FsController) render() {
	if c.len() == 0 {
		c.printEmptyFolder(c.Active())
		return
	}

	if c.IsChunk() {
		c.fullRender()
	} else {

		if c.PrevIndex() >= c.ChunkStart() && c.PrevIndex() <= c.ChunkEnd() {

			prevOpts := UpdateLineOptions{
				en:                 *c.get(c.PrevIndex()),
				index:              c.PrevIndex(),
				isActive:           false,
				isControllerActive: c.Active(),
				isSelected:         c.selected.IsSelected((*c.get(c.PrevIndex())).FullPath()),
				isLast:             c.PrevIndex() == c.len()-1,
			}

			c.RenderSingle(prevOpts)
		}

		opts := UpdateLineOptions{
			en:                 *c.current(),
			index:              c.Index(),
			isActive:           true,
			isControllerActive: c.Active(),
			isSelected:         c.selected.IsSelected((*c.current()).FullPath()),
			isLast:             c.Index() == c.len()-1,
		}

		c.RenderSingle(opts)

	}
}
