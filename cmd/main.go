package main

import (
	"fmt"
	"os"
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
	// navigators           []navigator.Navigator
	leftNav, rightNav, nav *navigator.Navigator
	// nav, inactiveNav             *navigator.Navigator

	done = make(chan bool)
)

func init() {
	c.ParseConfigFile(c.Cfg)
	c.ParseColorSchema(c.Cfg.CurrentSchema, &terminal.CurrentTheme)

	paneWidth = terminal.GetPaneWidth()
	//	nav = *navigator.NewNavigator()
	leftNav = navigator.NewNavigator()
	leftNav.StartCell = 2
	leftNav.RowWidth = paneWidth - 2

	nav = leftNav
	nav.IsActive = true
	// leftNav.IsActive = true
	// nav.StartCell = 2
	// nav.RowWidth = paneWidth - 2
	rightNav = navigator.NewNavigator()
	rightNav.StartCell = paneWidth + 2
	rightNav.RowWidth = paneWidth - 2

	// navigators = append(navigators, *navigator.NewNavigator())
	// navigators[0].StartCell = 2
	// navigators[0].RowWidth = paneWidth - 2
	// nav = &navigators[0]
	// nav.IsActive = true
	// navigators = append(navigators, *navigator.NewNavigator())
	// navigators[1].StartCell = paneWidth + 2
	// navigators[1].RowWidth = paneWidth - 2

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

func scan(n *navigator.Navigator, s *navigator.Selected) error {
	newRootDir, newEntries, err := task.ScanInputDirectory(c.Cfg.DefaultRootDirectory)
	if err != nil {
		return err
	}

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

	return nil
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

func showErrorAndTerminate(err error) {
	terminal.ClearScreen()
	fmt.Println(err)

	exit <- true
}

func navigate() {
	defer func() {
		if r := recover(); r != nil {
			terminal.ClearScreen()
			fmt.Println("Panic recovered ", r)
			return

		}
	}()

	keysEvents, err := k.GetKeys(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		// terminal.ShowCursor()

		err := k.Close()
		//	fmt.Print("Closing keyboard")
		if err != nil {
			fmt.Println("Error closing keyboard")
		}
	}()

	for keyListener {
		//	task.StartTicker()

		select {
		case <-pauseNavigation:
			keyListener = false
			//	fmt.Print("Pausing navigation")
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

				// navigators[1].Entries = append(navigators[1].Entries, nav.Entries...)
				// navigators[1].DirSize = nav.DirSize
				// navigators[1].CurrentPath = nav.CurrentPath
				// navigators[1].RootPath = nav.RootPath
				rightNav.Entries = append(rightNav.Entries, nav.Entries...)
				rightNav.DirSize = nav.DirSize
				rightNav.CurrentPath = nav.CurrentPath
				rightNav.RootPath = nav.RootPath

				terminal.RenderOutput(rightNav, &selected)
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

				// if activePane == 0 {
				// 	activePane = 1
				// 	nav.IsActive = false
				// 	terminal.RenderOutput(nav, &selected)
				// 	navigators[0] = *nav
				// 	*nav = navigators[1]
				// 	nav.IsActive = true

				// } else {
				// 	activePane = 0
				// 	nav.IsActive = false
				// 	terminal.RenderOutput(nav, &selected)
				// 	navigators[1] = *nav

				// 	*nav = navigators[0]
				// 	nav.IsActive = true
				// }

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
				pause()
				go func() {
					err := scan(nav, &selected)
					if err != nil {
						panic(err)
					}
					resume()
				}()

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
							}
							resume()
						}()
					}
				}

			case renameKey:
				if nav.HasEntries() {

					pause()
					go func() {
						ok, err := task.Rename(nav.GetCurrentEntry().Name, nav.CurrentPath)
						if err != nil {
							panic(err)
						}

						if ok {
							refresh(nav, &selected)
						}

						resume()
					}()

				}

			case enterKey:
				if nav.HasEntries() {
					switch nav.GetCurrentEntry().IsDir {
					case false:
						pause()
						go func() {
							task.ExecuteDefault(nav.GetCurrentEntry().FullPath())
							resume()
						}()

					default:

						enterSubfolder(nav, &selected)
					}
				}
			case deleteKey, deleteKey2, deleteKey3:

				selected.DumpPrevious(nav.CurrentPath)
				if !selected.IsEmpty() {
					pause()
					go func() {
						ok, err := task.DeleteSelected(&selected, nav)
						if err != nil {
							panic(err)
						}

						if ok {
							nav.CurrentIndex = 0
							refresh(nav, &selected)
						}

						//	terminal.ResetFlushOutput(nav, &selected)
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

				terminal.PrintPane(2, 1, paneWidth)
				terminal.PrintPane(2, paneWidth+1, paneWidth*2)

				terminal.RenderOutput(leftNav, &selected)
				terminal.RenderOutput(rightNav, &selected)

			case execKey:
				pause()

				go func() {
					input, err := task.WaitInput("run command: ", "")
					if err != nil {
						panic(err)
					}

					terminal.ClearScreen()
					task.ExecCommand(input)

					resume()
				}()

				// terminal.ResetFlushOutput(nav, &selected)

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
					ok, err := task.Relocate(prompt, rem, &selected, nav.CurrentPath)
					if err != nil {
						panic(err)
					}

					if ok {
						refresh(nav, &selected)
					}

					resume()
				}()

			case newFileKey:
				pause()
				go func() {
					ok, err := task.CreateFsFile(nav.CurrentPath)
					if err != nil {
						panic(err)
					}

					if ok {
						refresh(nav, &selected)
					}

					resume()
				}()

			case newDirKey:

				pause()
				go func() {
					ok, err := task.CreateFsDirectory(nav.CurrentPath)
					if err != nil {
						panic(err)
					}

					if ok {
						refresh(nav, &selected)
					}

					resume()
				}()

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
				terminal.ClearScreen()
				keyListener = true
				terminal.PrintPane(2, 1, paneWidth)
				terminal.PrintPane(2, paneWidth+1, paneWidth*2)
				terminal.RenderOutput(leftNav, &selected)
				terminal.RenderOutput(rightNav, &selected)

				navigate()
			}()
		}
	}
}
