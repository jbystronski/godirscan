package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/cache"
	"github.com/jbystronski/godirscan/pkg/config"
	c "github.com/jbystronski/godirscan/pkg/config"
	"github.com/jbystronski/godirscan/pkg/converter"
	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/navigator"
	"github.com/jbystronski/godirscan/pkg/task"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

var (
	execKey   = k.KeyCtrlE
	deleteKey = k.KeyDelete

	deleteKey3 = k.KeyCtrlD
	acceptKey  = k.KeyEnter
	enterKey   = acceptKey
	rejectKey  = k.KeyEsc

	renameKey = k.KeyCtrlR
	quitKey   = rejectKey

	spaceKey   = k.KeySpace
	newFileKey = k.KeyCtrlW

	findKey     = k.KeyCtrlF
	findSizeKey = k.KeyCtrlL
	sortKey     = k.KeyCtrlS

	menuKey = k.KeyF2

	editKey = k.KeyCtrl4
	copyKey = k.KeyCtrlV

	previewKey = k.KeyCtrlQ

	moveKey    = k.KeyF6
	newDirKey  = k.KeyF7
	deleteKey2 = k.KeyF8
	scanKey    = k.KeyF9
	quitKey2   = k.KeyF10

	homeKey       = k.KeyHome
	endKey        = k.KeyEnd
	pgDownKey     = k.KeyPgdn
	pgUpKey       = k.KeyPgup
	downKey       = k.KeyArrowDown
	leftKey       = k.KeyArrowLeft
	upKey         = k.KeyArrowUp
	rightKey      = k.KeyArrowRight
	selectKey     = k.KeyInsert
	selectKey2    = k.KeyCtrlI
	selectAllKey  = k.KeyCtrlA
	matchKey      = k.KeyCtrlO
	backSpaceKey  = k.KeyBackspace
	backSpaceKey2 = k.KeyBackspace2
	nextThemeKey  = k.KeyCtrlSlash
	switchPaneKey = k.KeyTab

	rejectChar = 'n'
	acceptChar = 'y'
	quitChar   = 'q'
)

var (
	firstRender         = true
	selected            navigator.Selected
	wg                  sync.WaitGroup
	sizeCalculationDone = make(chan struct{})
	searchDone          = make(chan struct{})
	recoveryChan        = make(chan struct{})

	resizeChan = make(chan os.Signal, 1)

	leftNav, rightNav, nav *navigator.Navigator
	done                   = make(chan bool)
)

func init() {
	c.ParseConfigFile(c.Cfg)
	c.ParseColorSchema(c.Cfg.CurrentSchema, &terminal.CurrentTheme)
	terminal.SetLayout()
	// paneWidth = terminal.GetPaneWidth()

	leftNav = navigator.NewNavigator()
	leftNav.StartCell = 2
	leftNav.RowWidth = terminal.GetPaneWidth() - 2

	nav = leftNav
	nav.IsActive = true

	rightNav = navigator.NewNavigator()
	rightNav.StartCell = terminal.GetPaneWidth() + 2
	rightNav.RowWidth = terminal.GetPaneWidth() - 2

	selected = *navigator.NewSelected()
}

func enterSubfolder(nav *navigator.Navigator, selected *navigator.Selected) {
	if nav.HasEntries() && nav.GetCurrentEntry().IsDir {

		p := nav.GetCurrentEntry().FullPath()

		if cachedEntries, ok := cache.Get(p); ok {
			nav.Entries = cachedEntries.Entries
			nav.DirSize = cachedEntries.Size
			nav.CurrentPath = p
			nav.AddBackTrace(nav.CurrentIndex)
			nav.CurrentIndex = 0
			terminal.ResetFlushOutput(nav, selected)
		} else {

			newPath, newEntries, err := task.ScanDirectory(p)
			if err != nil {
				terminal.FlashError(err)
				return
			}

			nav.CurrentPath = newPath
			nav.AddBackTrace(nav.CurrentIndex)

			nav.Entries = newEntries
			nav.SortMode = 0
			nav.CurrentIndex = 0
			entry.SetSort(&nav.SortMode, nav.Entries)

			terminal.RenderOutput(nav, selected)
			task.StartTicker()
			go func() {
				task.ScanDirectorySize(nav.Entries, &nav.DirSize)
				sizeCalculationDone <- struct{}{}
			}()

		}

	}
}

func execute(n *navigator.Navigator) {
	for {
		switch n.GetCurrentEntry().IsDir {
		case false:

			task.ExecuteDefault(n.GetCurrentEntry().FullPath())
			rerender()
			continue

		default:

			enterSubfolder(nav, &selected)
		}
	}
}

func scan(n *navigator.Navigator, s *navigator.Selected, offsetRow, offsetCol int) (string, error) {
	newRootDir, newEntries, err := task.ScanInputDirectory(c.Cfg.DefaultRootDirectory, offsetRow, offsetCol)
	if err != nil {
		return "", err
	}

	if newRootDir != "" {
		s.Clear()
		cache.Clear()
		n.CurrentPath = newRootDir
		n.RootPath = newRootDir
		n.Entries = newEntries
		entry.SetSort(&nav.SortMode, nav.Entries)
		rerender()
		//	terminal.ResetFlushOutput(n, s)
		task.StartTicker()
		go func() {
			task.ScanDirectorySize(n.Entries, &n.DirSize)

			sizeCalculationDone <- struct{}{}
		}()
	}

	return newRootDir, nil
}

func navigate() {
	defer func() {
		if r := recover(); r != nil {
			terminal.FlashError(fmt.Errorf("%s", r))
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
			terminal.FlashError(err)
		}
	}()

	signal.Notify(resizeChan, syscall.SIGWINCH)

	for {
		select {

		case <-recoveryChan:
			fmt.Println("rerendirg after recovery")
			time.Sleep(time.Second * 2)
			rerender()

		case <-resizeChan:

			terminal.SetLayout()
			rerender()

		case <-sizeCalculationDone:

			task.StopTicker()
			//	cache.Store(nav.CurrentPath, *&nav.DirSize, nav.Entries)
			terminal.RenderOutput(nav, &selected)
			if firstRender {

				rightNav.Entries = append(rightNav.Entries, nav.Entries...)
				rightNav.DirSize = nav.DirSize
				rightNav.CurrentPath = nav.CurrentPath
				rightNav.RootPath = nav.RootPath

				//	terminal.RenderOutput(rightNav, &selected)
				rerender()
				firstRender = false
			}

		case <-task.Ticker.C:
			rerender()
			// terminal.RenderOutput(nav, &selected)

		case <-searchDone:
			if nav.GetEntriesLength() == 0 {
				rerender()
				fmt.Println("no entries found")
			}

		case event := <-keysEvents:

			switch event.Key {

			case switchPaneKey:

				nav.IsActive = false

				if nav == leftNav {
					leftNav.IsActive = false
					nav = rightNav
				} else {
					rightNav.IsActive = false
					nav = leftNav
				}
				nav.IsActive = true

				terminal.RenderOutput(leftNav, &selected)

				terminal.RenderOutput(rightNav, &selected)
			case menuKey:

				terminal.ClearScreen()
				terminal.PrintHelp()

				rerender()

			case k.KeyCtrlC:
				terminal.ClearScreen()
				return
			case scanKey:

				_, err := scan(nav, &selected, terminal.PromptLine, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					return
				}
				refresh(nav, &selected)

			case selectKey:
				if nav.HasEntries() {
					selected.DumpPrevious(nav.CurrentPath)
					selected.Select(nav.GetCurrentEntry())
					nav.MoveDown()
					terminal.RenderOutput(nav, &selected)

				}

			case selectAllKey:
				if nav.HasEntries() {
					selected.DumpPrevious(nav.CurrentPath)
					selected.SelectAll(nav.Entries)
					terminal.RenderOutput(nav, &selected)
				}

			case downKey:
				if nav.MoveDown() {
					terminal.RenderOutput(nav, &selected)
				}

			case upKey:
				if nav.MoveUp() {
					terminal.RenderOutput(nav, &selected)
				}
			case rightKey:
				enterSubfolder(nav, &selected)

			case leftKey:

				if cachedEntries, ok := cache.Get(nav.GetParentPath()); ok {
					nav.Entries = cachedEntries.Entries
					nav.DirSize = cachedEntries.Size
					nav.CurrentPath = nav.GetParentPath()

					nav.CurrentIndex = nav.GetBackTrace()

					nav.SortMode = 0
					entry.SetSort(&nav.SortMode, nav.Entries)
					terminal.RenderOutput(nav, &selected)

				} else {

					newPath, newEntries, err := task.ScanDirectory(nav.GetParentPath())
					if err != nil {
						terminal.FlashError(err)
						time.Sleep(time.Second * 2)
						return
					}
					nav.CurrentPath = newPath
					nav.Entries = newEntries

					nav.CurrentIndex = nav.GetBackTrace()

					nav.SortMode = 0
					entry.SetSort(&nav.SortMode, nav.Entries)
					terminal.RenderOutput(nav, &selected)
					task.StartTicker()
					go func() {
						task.ScanDirectorySize(nav.Entries, &nav.DirSize)
						sizeCalculationDone <- struct{}{}
					}()

				}

			case findKey:

				path, err := task.WaitInput("find in path", nav.CurrentPath, terminal.PromptLine, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					continue
				}
				if path == "" {
					continue
				}
				pattern, err := task.WaitInput("find (pattern)", "", terminal.PromptLine, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					continue
				}
				if pattern == "" {
					continue
				}

				var find func(*regexp.Regexp, string, *[]*entry.Entry) error

				find = func(reg *regexp.Regexp, path string, entries *[]*entry.Entry) error {
					dc, err := os.ReadDir(path)
					if err != nil {
						return err
					}

					for _, en := range dc {

						info, err := en.Info()
						if err != nil {
							terminal.FlashError(err)
							continue
						}

						if info.IsDir() {
							find(reg, filepath.Join(path, en.Name()), entries)
						} else {
							if reg.Match([]byte(info.Name())) {
								*entries = append(*entries, &entry.Entry{
									Name:  filepath.Join(path, info.Name()),
									Size:  int(info.Size()),
									IsDir: info.IsDir(),
									Path:  &path,
								})
							}
						}
					}
					return nil
				}
				compiled := regexp.MustCompile(pattern)
				nav.Entries = nil
				nav.CurrentIndex = 0

				terminal.RenderOutput(nav, &selected)

				task.StartTicker()

				go func() {
					err := find(compiled, filepath.Join(path), &nav.Entries)
					if err != nil {
						terminal.FlashError(err)
						return
					}

					searchDone <- struct{}{}
				}()

			case findSizeKey:

				answ, err := task.WaitInput("Find by size, unit: ( 0=bytes 1=kb 2=mb 3=gb ) ", "2", terminal.PromptLine, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				if answ == "" {
					continue
				}

				unit, _ := strconv.Atoi(answ)

				if unit < 0 || unit > len(converter.StorageUnits)-1 {
					terminal.FlashError(errors.New("invalid unit index"))
					time.Sleep(time.Second * 2)

					continue
				}

				unitName := converter.StorageUnits[unit]

				answ, err = task.WaitInput(fmt.Sprintf("%s %s", "Type min value in", unitName), "0", terminal.PromptLine, nav.StartCell)

				if answ == "" {
					continue
				}

				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				min, err := strconv.ParseFloat(answ, 64)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				if min < 0 {
					min = 0
				}

				answ, err = task.WaitInput(fmt.Sprintf("%s %s", "Type max value ( 0 or no value means unlimited ) in ", unitName), "", terminal.PromptLine, nav.StartCell)

				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				if answ == "" {
					answ = "0"
				}

				max, err := strconv.ParseFloat(answ, 64)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				if max < min {
					terminal.FlashError(fmt.Errorf("Max value: %v, can't be lower than min value: %v", max, min))
					time.Sleep(time.Second * 2)
					continue
				}

				path, err := task.WaitInput("Directory to search from: ", "", terminal.PromptLine, nav.StartCell)

				if path == "" {
					continue
				}

				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				pattern, err := task.WaitInput("Pattern to match: ", "", terminal.PromptLine, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
					continue
				}

				compiled := regexp.MustCompile(pattern)

				minV, maxV := converter.ToBytes(converter.StorageUnits[unit], min, max)

				var find func(string, *regexp.Regexp, int64, int64, *[]*entry.Entry) error

				find = func(path string, reg *regexp.Regexp, min, max int64, entries *[]*entry.Entry) error {
					dc, err := os.ReadDir(path)
					if err != nil {
						return err
					}

					for _, dirEntry := range dc {

						info, _ := dirEntry.Info()

						if info.IsDir() {
							find(filepath.Join(path, info.Name()), reg, min, max, entries)
						} else {
							if info.Size() >= min {
								if max != 0 && info.Size() <= max || max == 0 {
									if reg.Match([]byte(info.Name())) {

										newEntry := &entry.Entry{
											Name:  filepath.Join(path, info.Name()),
											Size:  int(info.Size()),
											IsDir: dirEntry.IsDir(),
											Path:  &path,
										}

										*entries = append(*entries, newEntry)

									}
								}
							}
						}
					}
					return nil
				}

				nav.Entries = nil
				nav.CurrentIndex = 0

				task.StartTicker()
				go func() {
					err := find(filepath.Join(path), compiled, minV, maxV, &nav.Entries)
					if err != nil {
						panic(err)
					}
					searchDone <- struct{}{}
				}()

			case sortKey:
				if nav.HasEntries() {

					entry.SetSort(&nav.SortMode, nav.Entries)
					terminal.RenderOutput(nav, &selected)

				}

			case homeKey:
				if nav.HasEntries() {
					nav.CurrentIndex = 0

					terminal.RenderOutput(nav, &selected)
				}

			case endKey:
				if nav.HasEntries() {
					nav.CurrentIndex = nav.GetEntriesLength() - 1
					terminal.RenderOutput(nav, &selected)

				}

			case pgDownKey:
				if nav.HasEntries() {
					if nav.CurrentIndex+terminal.OutputLines >= nav.GetEntriesLength() {
						nav.CurrentIndex = nav.GetEntriesLength() - 1
					} else {
						nav.CurrentIndex += terminal.OutputLines
					}
					terminal.RenderOutput(nav, &selected)
				}

			case pgUpKey:
				if nav.HasEntries() {
					if nav.CurrentIndex-terminal.OutputLines < 0 {
						nav.CurrentIndex = 0
					} else {
						nav.CurrentIndex -= terminal.OutputLines
					}
					terminal.RenderOutput(nav, &selected)

				}

			case matchKey:
				if nav.HasEntries() {
					var matches []int
					var curr int

					test, err := task.WaitInput("Match: ", "", terminal.PromptLine, nav.StartCell)
					if err != nil || test == "" {
						terminal.FlashError(err)
						continue
					}

					for index, en := range nav.Entries {
						if strings.Contains(en.Name, test) {
							matches = append(matches, index)
						}
					}

					fmt.Print(len(matches), " navigate")

					if len(matches) > 0 {
						nav.CurrentIndex = matches[0]
						terminal.RenderOutput(nav, &selected)
					} else {
						terminal.FlashError(errors.New("no matches found"))
						terminal.ClearLine()
						terminal.CarriageReturn()
						continue
					}

				MatchLoop:
					for {

						_, key, err := k.GetKey()
						if err != nil {
							terminal.FlashError(err)
							break
						}
						switch key {
						case upKey:
							if curr == 0 {
								continue
							}

							curr--
							nav.CurrentIndex = matches[curr]
							terminal.RenderOutput(nav, &selected)
						case downKey:
							if curr == len(matches)-1 {
								continue
							}

							curr++
							nav.CurrentIndex = matches[curr]
							terminal.RenderOutput(nav, &selected)

						case k.KeyEsc, enterKey:
							terminal.ClearLine()
							terminal.CarriageReturn()
							break MatchLoop
						}

					}

					continue
				}

			case editKey:

				if nav.HasEntries() {
					if !nav.GetCurrentEntry().IsDir {

						terminal.ClearScreen()
						sizeBefore := nav.GetCurrentEntry().Size
						task.Edit(nav.GetCurrentEntry().FullPath(), c.Cfg.DefaultEditor)

						info, _ := os.Stat(nav.GetCurrentEntry().FullPath())

						if info.Size() != int64(sizeBefore) {
							refresh(nav, &selected)
						} else {
							rerender()
						}
						continue

					}
				}

			case renameKey:
				if nav.HasEntries() {

					ok, err := task.Rename(nav.GetCurrentEntry().Name, nav.CurrentPath, nav.StartCell)
					if err != nil {
						terminal.FlashError(err)
						time.Sleep(time.Second * 2)
					}

					if ok {
						refresh(nav, &selected)
					}

				}

			case enterKey:
				if nav.HasEntries() {
					// executeDefault(nav)
					switch nav.GetCurrentEntry().IsDir {
					case false:

						task.ExecuteDefault(nav.GetCurrentEntry().FullPath())
						rerender()
						continue

					default:

						enterSubfolder(nav, &selected)
					}
				}
			case deleteKey, deleteKey2, deleteKey3:

				selected.DumpPrevious(nav.CurrentPath)
				if !selected.IsEmpty() {
					{
						ok, err := task.DeleteSelected(&selected, nav)
						if err != nil {
							terminal.FlashError(err)
							time.Sleep(time.Second * 2)
						}

						if ok {
							nav.CurrentIndex = 0
							refresh(nav, &selected)
						}

					}
				}

			case nextThemeKey:
				num := config.Cfg.CurrentSchema
				if num < uint(len(config.Cfg.ColorSchemas)-1) {
					num++
				} else {
					num = 0
				}

				config.Cfg.CurrentSchema = num
				config.ParseColorSchema(num, &terminal.CurrentTheme)
				config.UpdateConfigFile(config.Cfg)
				rerender()

			case execKey:

				input, err := task.WaitInput("run command: ", "", terminal.PromptLine, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					//	time.Sleep(time.Second * 2)
				} else {
					terminal.ClearScreen()
					task.ExecCommand(input)
					rerender()
					continue

				}

			case copyKey, moveKey:

				if selected.IsEmpty() {
					continue
				}

				var prompt string
				var rem bool

				if event.Key == moveKey {
					prompt = "Move"
					rem = true
				} else {
					prompt = "Copy"
				}

				ok, err := task.Relocate(prompt, rem, &selected, nav.CurrentPath, nav.StartCell)
				if err != nil {
					panic(err)
					// fmt.Println("ERROR OCCURED")
					// time.Sleep(time.Second * 3)
					// terminal.FlashError(err)
					// time.Sleep(time.Second * 3)
				}

				if ok {
					refresh(nav, &selected)
				}

			case newFileKey:

				ok, err := task.CreateFsFile(nav.CurrentPath, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
				}

				if ok {
					refresh(nav, &selected)
				}

			case newDirKey:

				ok, err := task.CreateFsDirectory(nav.CurrentPath, nav.StartCell)
				if err != nil {
					terminal.FlashError(err)
					time.Sleep(time.Second * 2)
				}

				if ok {
					refresh(nav, &selected)
				}

			case quitKey, quitKey2:
				terminal.ClearScreen()
				return

			}

		}
	}
}

func main() {
	defer func() {
		if k.IsStarted(time.Millisecond * 50) {
			k.Close()
		}
	}()

	printStart()
	navigate()
}

func rerender() {
	terminal.ClearScreen()
	leftNav.StartCell = 2
	leftNav.RowWidth = terminal.GetPaneWidth() - 2
	rightNav.StartCell = terminal.GetPaneWidth() + 2
	rightNav.RowWidth = terminal.GetPaneWidth() - 2

	terminal.PrintPanes()
	terminal.RenderOutput(leftNav, &selected)
	terminal.RenderOutput(rightNav, &selected)
}

func refresh(nav *navigator.Navigator, selected *navigator.Selected) {
	var err error

	nav.CurrentPath, nav.Entries, err = task.ScanDirectory(nav.CurrentPath)
	if err != nil {
		terminal.FlashError(err)
		time.Sleep(time.Second * 2)
		return
	}
	selected.Clear()
	cache.Clear()

	nav.SortMode = 0

	entry.SetSort(&nav.SortMode, nav.Entries)
	rerender()
	task.StartTicker()
	go func() {
		task.ScanDirectorySize(nav.Entries, &nav.DirSize)
		sizeCalculationDone <- struct{}{}
	}()
}

func printStart() {
	_ = k.Open()

	defer func() {
		k.Close()
	}()

	terminal.ClearScreen()

	terminal.PrintBanner()

	root, err := scan(nav, &selected, 10, 1)

	if root == "" || err != nil {
		terminal.ClearScreen()
		k.Close()
		os.Exit(0)
	}
	terminal.ClearScreen()
	terminal.PrintPanes()
}
