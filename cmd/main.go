package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/cache"
	"github.com/jbystronski/godirscan/pkg/config"
	c "github.com/jbystronski/godirscan/pkg/config"
	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/navigator"
	"github.com/jbystronski/godirscan/pkg/task"
	"github.com/jbystronski/godirscan/pkg/terminal"
	"github.com/jbystronski/godirscan/pkg/utils"
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
	viewKey = k.KeyF3
	editKey = k.KeyCtrl4
	copyKey = k.KeyCtrlV

	previewKey = k.KeyCtrlQ

	moveKey       = k.KeyF6
	newDirKey     = k.KeyF7
	deleteKey2    = k.KeyF8
	scanKey       = k.KeyF9
	quitKey2      = k.KeyF10
	goToKey       = k.KeyCtrlG
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
	backSpaceKey  = k.KeyBackspace
	backSpaceKey2 = k.KeyBackspace2
	nextThemeKey  = k.KeyCtrlSlash
	switchPaneKey = k.KeyTab

	rejectChar = 'n'
	acceptChar = 'y'
	quitChar   = 'q'
)

var (
	// nav                  navigator.Navigator
	firstRender          = true
	selected             navigator.Selected
	wg                   sync.WaitGroup
	sizeCalculationDone  = make(chan struct{})
	searchDone           = make(chan struct{})
	pauseNavigation      = make(chan struct{}, 1)
	stopNavigation       = make(chan struct{}, 1)
	resumeNavigationChan = make(chan struct{}, 1)
	taskDone             = make(chan func())
	exit                 = make(chan bool, 1)
	paneWidth            int
	activePane           = 0
	navigators           []navigator.Navigator
	nav                  *navigator.Navigator

	done = make(chan bool)
)

func init() {
	c.ParseConfigFile(c.Cfg)
	c.ParseColorSchema(c.Cfg.CurrentSchema, &terminal.CurrentTheme)

	paneWidth = terminal.GetPaneWidth()
	//	nav = *navigator.NewNavigator()
	navigators = append(navigators, *navigator.NewNavigator())
	navigators[0].StartCell = 2
	navigators[0].RowWidth = paneWidth - 2
	nav = &navigators[0]
	nav.IsActive = true
	navigators = append(navigators, *navigator.NewNavigator())
	navigators[1].StartCell = paneWidth + 2
	navigators[1].RowWidth = paneWidth - 2

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
			terminal.ResetFlushOutput(nav, selected)
		} else {
			task.StartTicker()
			newPath, newEntries := task.ScanDirectory(p)

			nav.CurrentPath = newPath
			nav.Entries = newEntries
			nav.SortMode = 0

			entry.SetSort(&nav.SortMode, nav.Entries)
			nav.AddBackTrace(nav.CurrentIndex)
			terminal.ResetFlushOutput(nav, selected)

			go func() {
				task.ScanDirectorySize(nav.Entries, &nav.DirSize)
				sizeCalculationDone <- struct{}{}
			}()

		}

	}
}

func scan(n *navigator.Navigator, s *navigator.Selected) {
	newRootDir, newEntries := task.ScanInputDirectory(c.Cfg.DefaultRootDirectory)

	if newRootDir != "" {
		s.Clear()
		cache.Clear()
		n.CurrentPath = newRootDir
		n.RootPath = newRootDir
		n.Entries = newEntries
		entry.SetSort(&nav.SortMode, nav.Entries)

		terminal.ResetFlushOutput(n, s)
		task.StartTicker()
		go func() {
			task.ScanDirectorySize(n.Entries, &n.DirSize)

			sizeCalculationDone <- struct{}{}
		}()
	}
}

func refresh(nav *navigator.Navigator, selected *navigator.Selected) {
	selected.Clear()
	cache.Clear()
	task.StartTicker()
	nav.CurrentPath, nav.Entries = task.ScanDirectory(nav.CurrentPath)

	go func() {
		task.ScanDirectorySize(nav.Entries, &nav.DirSize)
		sizeCalculationDone <- struct{}{}
	}()
}

var keyListener = true

func navigate() {
	keysEvents, err := k.GetKeys(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		// terminal.ShowCursor()

		err := k.Close()
		if err != nil {
			fmt.Println("Error closing keyboard")
		}
	}()

	for keyListener {
		//	task.StartTicker()

		select {
		case <-pauseNavigation:
			keyListener = false
			break
		case <-stopNavigation:
			keyListener = false
			exit <- true
			break
		case <-sizeCalculationDone:

			task.StopTicker()
			cache.Store(nav.CurrentPath, *&nav.DirSize, nav.Entries)
			terminal.RenderOutput(nav, &selected)
			if firstRender {

				navigators[1].Entries = append(navigators[1].Entries, nav.Entries...)
				navigators[1].DirSize = nav.DirSize
				navigators[1].CurrentPath = nav.CurrentPath
				navigators[1].RootPath = nav.RootPath

				terminal.RenderOutput(&navigators[1], &selected)
				firstRender = false
			}

		case <-task.Ticker.C:

			terminal.RenderOutput(nav, &selected)

			// if firstRender {
			// 	terminal.RenderOutput(&navigators[1], &selected)
			// }

		case <-searchDone:
			if nav.GetEntriesLength() == 0 {
				fmt.Println("no entries found")
			}

		case event := <-keysEvents:
			switch event.Key {

			case switchPaneKey:

				if activePane == 0 {
					activePane = 1
					nav.IsActive = false
					terminal.RenderOutput(nav, &selected)
					navigators[0] = *nav
					*nav = navigators[1]
					nav.IsActive = true

				} else {
					activePane = 0
					nav.IsActive = false
					terminal.RenderOutput(nav, &selected)
					navigators[1] = *nav

					*nav = navigators[0]
					nav.IsActive = true
				}

				terminal.RenderOutput(nav, &selected)
			case menuKey:

				go func() {
					terminal.ClearScreen()
					terminal.PrintHelp()

					resumeNavigation()
					terminal.ResetFlushOutput(nav, &selected)
				}()
				pauseNavigation <- struct{}{}

			case k.KeyCtrlC:
				return
			case scanKey:
				go func() {
					scan(nav, &selected)
					resume()
				}()

				pause()
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

				if nav.CurrentPath != nav.RootPath {
					if cachedEntries, ok := cache.Get(nav.GetParentPath()); ok {
						nav.Entries = cachedEntries.Entries
						nav.DirSize = cachedEntries.Size
						nav.CurrentPath = nav.GetParentPath()
						nav.StartLine = 0
						nav.EndLine = 0
						nav.CurrentIndex = nav.GetBackTrace()

						nav.SortMode = 0
						entry.SetSort(&nav.SortMode, nav.Entries)
						terminal.FlushOutput(nav, &selected)

					} else {
						task.StartTicker()
						newPath, newEntries := task.ScanDirectory(nav.GetParentPath())
						nav.CurrentPath = newPath
						nav.Entries = newEntries
						nav.StartLine = 0
						nav.EndLine = 0
						nav.CurrentIndex = nav.GetBackTrace()

						nav.SortMode = 0
						entry.SetSort(&nav.SortMode, nav.Entries)
						terminal.FlushOutput(nav, &selected)

						go func() {
							task.ScanDirectorySize(nav.Entries, &nav.DirSize)
							sizeCalculationDone <- struct{}{}
						}()

					}
				}

			case findKey:
				task.PromptFindByName(nav.CurrentPath, nav.Entries, searchDone)

			// case findSizeKey:

			// 	task.PromptFindBySize()

			// 	go func() {
			// 		task.FindBySize(nav.Entries, pathName, pattern, minV, maxV)
			// 		searchDone <- struct{}{}
			// 	}()

			case sortKey:
				if nav.HasEntries() {

					entry.SetSort(&nav.SortMode, nav.Entries)
					terminal.RenderOutput(nav, &selected)

				}

			case homeKey:
				if nav.HasEntries() {
					nav.CurrentIndex = 0
					nav.StartLine = 0
					terminal.RenderOutput(nav, &selected)
				}

			case endKey:
				if nav.HasEntries() {
					nav.CurrentIndex = nav.GetEntriesLength() - 1
					if nav.GetEntriesLength() < nav.NumVisibleLines {
						nav.StartLine = 0
					} else {
						nav.StartLine = nav.GetEntriesLength() - nav.NumVisibleLines
					}
				}

			case pgDownKey:
				if nav.HasEntries() {
					terminal.ClearScreen()
					if nav.GetEntriesLength() < nav.CurrentIndex+nav.NumVisibleLines {
						nav.SetCurrentIndex(nav.GetEntriesLength() - 1)
						nav.SetStartLine(nav.GetEntriesLength() - nav.NumVisibleLines)
					} else {
						nav.SetCurrentIndex(nav.CurrentIndex + nav.NumVisibleLines)
						nav.SetStartLine(nav.CurrentIndex)
					}

					terminal.RenderOutput(nav, &selected)
				}

			case pgUpKey:
				if nav.HasEntries() {
					if nav.GetCurrentIndex()-nav.NumVisibleLines < 0 {
						nav.SetCurrentIndex(0)
						nav.SetStartLine(0)
					} else {
						nav.SetCurrentIndex(nav.CurrentIndex - nav.NumVisibleLines)
						nav.SetStartLine(nav.CurrentIndex)
					}

					terminal.RenderOutput(nav, &selected)
				}

			case viewKey:
				if nav.HasEntries() {
					if !nav.GetCurrentEntry().IsDir {
						terminal.ClearScreen()
						task.Peek(nav.GetCurrentEntry().FullPath())
					}
				}

			case editKey:

				if nav.HasEntries() {
					if !nav.GetCurrentEntry().IsDir {
						pause()
						go func() {
							terminal.ClearScreen()
							sizeBefore := nav.GetCurrentEntry().Size
							task.Edit(nav.GetCurrentEntry().FullPath(), c.Cfg.DefaultEditor)

							info, _ := os.Stat(nav.GetCurrentEntry().FullPath())

							if info.Size() != int64(sizeBefore) {
								refresh(nav, &selected)
								terminal.ResetFlushOutput(nav, &selected)
							} else {
								terminal.RenderOutput(nav, &selected)
							}
							resume()
						}()
					}
				}

			case renameKey:
				if nav.HasEntries() {

					answ := task.WaitInput("rename", nav.GetCurrentEntry().Name)

					if strings.Contains(answ, string(os.PathSeparator)) {
						utils.ShowErrAndContinue(fmt.Errorf("path separator can't be used inside name"))

						terminal.ResetFlushOutput(nav, &selected)
						return
					}

					if answ == "" {
						utils.ShowErrAndContinue(fmt.Errorf("empty name?"))
						terminal.ResetFlushOutput(nav, &selected)
						return
					}

					_, err := os.Stat(filepath.Join(nav.CurrentPath, answ))
					if err == nil {

						utils.ShowErrAndContinue(fmt.Errorf("entry already exists"))
						terminal.ResetFlushOutput(nav, &selected)
						return

					}

					err = os.Rename(nav.GetCurrentEntry().FullPath(), filepath.Join(nav.CurrentPath, answ))

					if err != nil {
						panic(err)
					}

					nav.GetCurrentEntry().Name = answ
					terminal.ResetFlushOutput(nav, &selected)

				}

			case enterKey:
				if nav.HasEntries() {
					switch nav.GetCurrentEntry().IsDir {
					case false:
						task.ExecuteDefault(nav.GetCurrentEntry().FullPath())
					default:
						enterSubfolder(nav, &selected)
					}
				}
			case deleteKey, deleteKey2, deleteKey3:

				selected.DumpPrevious(nav.CurrentPath)
				if !selected.IsEmpty() {
					pause()
					go func() {
						task.DeleteSelected(&selected, nav)
						refresh(nav, &selected)
						terminal.ResetFlushOutput(nav, &selected)
						resume()
					}()
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

				terminal.RenderOutput(nav, &selected)

			case execKey:
				input := task.WaitInput("run command: ", "")

				terminal.ClearScreen()
				task.ExecCommand(input)

				terminal.ResetFlushOutput(nav, &selected)

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

				pause()

				go func() {
					task.Relocate(prompt, rem, &selected, nav)
					refresh(nav, &selected)
					resume()
				}()

			case newFileKey:

				go func() {
					if task.CreateFsFile(nav.CurrentPath) {
						//	sortMode = entry.SortByName(entries)

						refresh(nav, &selected)
					}
					resume()
				}()
				pause()
			case newDirKey:
				// keyListener = false
				// return
				pause()
				go func() {
					fmt.Println("Doing some intensive task..")
					time.Sleep(time.Second * 1)
					resumeNavigationChan <- struct{}{}
				}()

				// if task.CreateFsDirectory(nav.CurrentPath) {
				// 	refresh(nav, &selected)
				// }

			case quitKey, quitKey2:
				terminal.ClearScreen()
				stopNavigation <- struct{}{}
				// pauseNavigation <- struct{}{}

			}

		}
	}
}

func pause() {
	// fmt.Println("stopping navigation mode")
	// time.Sleep(time.Second * 1)
	pauseNavigation <- struct{}{}
	time.Sleep(time.Millisecond * 50)
}

func resume() {
	// fmt.Println("resuming navigation")
	// time.Sleep(time.Second * 1)
	resumeNavigationChan <- struct{}{}
}

func resumeNavigation() {
}

func main() {
	terminal.ClearScreen()

	terminal.PrintBanner()

	scan(nav, &selected)

	terminal.ClearScreen()

	terminal.PrintPane(2, 1, paneWidth)
	terminal.PrintPane(2, paneWidth+1, paneWidth*2)

	navigate()
	for {
		select {
		case <-exit:

			return
		case <-resumeNavigationChan:
			// fmt.Println("resume signal received")
			// time.Sleep(time.Second * 1)
			// wg.Add(1)

			go func() {
				// defer wg.Done()
				// fmt.Println("resuming navigation")
				// time.Sleep(time.Second * 1)
				// refresh(nav, &selected)

				keyListener = true
				// fmt.Println("keylistener on")
				// time.Sleep(time.Second * 1)
				navigate()
			}()
		}
	}
}
