package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/cache"
	"github.com/jbystronski/godirscan/pkg/config"
	c "github.com/jbystronski/godirscan/pkg/config"
	"github.com/jbystronski/godirscan/pkg/converter"
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
	editKey = k.KeyF4
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

	rejectChar = 'n'
	acceptChar = 'y'
	quitChar   = 'q'
)

var (
	nav navigator.Navigator

	selected            navigator.Selected
	sizeCalculationDone = make(chan struct{})
	searchDone          = make(chan struct{})
	stopListener        = make(chan struct{})
	done                = make(chan bool)
)

func init() {
	c.ParseConfigFile(c.Cfg)
	c.ParseColorSchema(c.Cfg.CurrentSchema, &terminal.CurrentTheme)
	nav = *navigator.NewNavigator()
	selected = *navigator.NewSelected()
}

func enterSubfolder(nav *navigator.Navigator, selected *navigator.Selected) {
	if nav.HasEntries() && nav.GetCurrentEntry().IsDir {

		p := nav.GetCurrentEntry().FullPath()

		if cachedEntries, ok := cache.Get(p); ok {
			nav.SetEntries(cachedEntries.Entries)
			nav.SetDirSize(cachedEntries.Size)
			nav.SetCurrentPath(p)
			nav.AddBackTrace(nav.GetCurrentIndex())
			terminal.ResetFlushOutput(nav, selected)
		} else {
			task.StartTicker()
			newPath, newEntries := task.ScanDirectory(p)

			nav.SetCurrentPath(newPath)
			nav.SetEntries(newEntries)
			nav.SortMode = 0

			entry.SetSort(&nav.SortMode, nav.GetEntries())
			nav.AddBackTrace(nav.GetCurrentIndex())
			terminal.ResetFlushOutput(nav, selected)

			go func() {
				task.ScanDirectorySize(nav.GetEntries(), nav.GetDirSize())
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
		n.SetCurrentPath(newRootDir)
		n.SetRootPath(newRootDir)
		n.SetEntries(newEntries)
		entry.SetSort(&nav.SortMode, nav.GetEntries())

		terminal.ResetFlushOutput(n, s)
		task.StartTicker()
		go func() {
			task.ScanDirectorySize(n.GetEntries(), n.GetDirSize())

			sizeCalculationDone <- struct{}{}
		}()
	}
}

var wg sync.WaitGroup

var keyListener = true

func main() {
	terminal.ClearScreen()

	refresh := func() {
		selected.Clear()
		cache.Clear()
		task.StartTicker()
		newPath, newEntries := task.ScanDirectory(nav.GetCurrentPath())
		nav.SetCurrentPath(newPath)
		nav.SetEntries(newEntries)

		go func() {
			task.ScanDirectorySize(nav.GetEntries(), nav.GetDirSize())
			sizeCalculationDone <- struct{}{}
		}()
	}

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

	terminal.PrintBanner()
	scan(&nav, &selected)

	for keyListener {

		task.StartTicker()

		select {
		case <-sizeCalculationDone:
			task.StopTicker()
			cache.Store(nav.GetCurrentPath(), *nav.GetDirSize(), nav.GetEntries())
			terminal.RenderOutput(&nav, &selected)
		case <-task.Ticker.C:
			terminal.RenderOutput(&nav, &selected)
		case <-searchDone:
			if nav.GetEntriesLength() == 0 {
				fmt.Println("no entries found")
			}

		case event := <-keysEvents:
			switch event.Key {
			case menuKey:
				terminal.ClearScreen()
				terminal.PrintHelp()
				terminal.ResetFlushOutput(&nav, &selected)
			case k.KeyCtrlC:
				return
			case scanKey:
				scan(&nav, &selected)

			case downKey:
				if nav.MoveDown() {
					terminal.RenderOutput(&nav, &selected)
				}

			case upKey:
				if nav.MoveUp() {
					terminal.RenderOutput(&nav, &selected)
				}
			case rightKey:
				enterSubfolder(&nav, &selected)

			case leftKey:
				if nav.GetCurrentPath() != nav.GetRootPath() {
					if cachedEntries, ok := cache.Get(nav.GetParentPath()); ok {
						nav.SetEntries(cachedEntries.Entries)
						nav.SetDirSize(cachedEntries.Size)
						nav.SetCurrentPath(nav.GetParentPath())
						nav.SetStartLine(0)
						nav.SetEndLine(0)
						nav.SetCurrentIndex(nav.GetBackTrace())
						nav.SortMode = 0
						entry.SetSort(&nav.SortMode, nav.GetEntries())
						terminal.FlushOutput(&nav, &selected)

					} else {
						task.StartTicker()
						newPath, newEntries := task.ScanDirectory(nav.GetParentPath())
						nav.SetCurrentPath(newPath)
						nav.SetEntries(newEntries)
						nav.SetStartLine(0)
						nav.SetEndLine(0)
						nav.SetCurrentIndex(nav.GetBackTrace())
						nav.SortMode = 0
						entry.SetSort(&nav.SortMode, nav.GetEntries())
						terminal.FlushOutput(&nav, &selected)

						go func() {
							task.ScanDirectorySize(nav.GetEntries(), nav.GetDirSize())
							sizeCalculationDone <- struct{}{}
						}()

					}
				}

			case findKey:
				task.PromptFindByName(nav.GetCurrentPath(), nav.GetEntries(), searchDone)

			case findSizeKey:

				var unitAsString, pathName, pattern string
				var unitAsInt int
				var min, max float64 = 0, math.MaxFloat64

				task.WaitUserInput("Find by size, unit: ( 0=bytes 1=kb 2=mb 3=gb ) ", "2", func(s string) {
					number, _ := strconv.Atoi(s)

					if number < 0 || number > len(converter.StorageUnits)-1 {
						utils.ShowErrAndContinue(errors.New("invalid unix index"))
						return
					}

					unitAsString = converter.StorageUnits[number]
					unitAsInt = number

					task.WaitUserInput(fmt.Sprintf("%s %s", "Type min value in", unitAsString), "0", func(s string) {
						num, err := strconv.ParseFloat(s, 64)
						if err != nil {
							utils.ShowErrAndContinue(err)
							return
						}

						if num < 0 {
							num = 0
						}
						min = num

						task.WaitUserInput(fmt.Sprintf("%s %s", "Type max value ( 0 or no value means unlimited ) in", unitAsString), "", func(s string) {
							if s != "" || s != "0" {
								n, err := strconv.ParseFloat(s, 64)
								if err != nil {
									utils.ShowErrAndContinue(err)
									return
								}
								max = n

							}

							if max < min {
								utils.ShowErrAndContinue(fmt.Errorf("max value: %v, can't be lower than min value: %v", max, min))
								return
							}

							task.WaitUserInput("Root directory to search from: ", "", func(s string) {
								_, err := os.Stat(s)
								if err != nil {
									utils.ShowErrAndContinue(err)
									return
								}

								pathName = s

								task.WaitUserInput("Pattern to match: ", "", func(s string) {
									pattern = s

									minV, maxV := converter.ToBytes(converter.StorageUnits[unitAsInt], min, max)

									task.StartTicker()

									nav.ClearEntries()
									terminal.ClearScreen()

									go func() {
										task.FindBySize(nav.GetEntries(), pathName, pattern, minV, maxV)
										searchDone <- struct{}{}
									}()
								})
							})
						})
					})
				})

			case sortKey:
				if nav.HasEntries() {

					entry.SetSort(&nav.SortMode, nav.GetEntries())
					terminal.RenderOutput(&nav, &selected)

				}

			case homeKey:
				if nav.HasEntries() {
					nav.SetCurrentIndex(0)
					nav.SetStartLine(0)
					terminal.RenderOutput(&nav, &selected)
				}

			case endKey:
				if nav.HasEntries() {
					nav.SetCurrentIndex(nav.GetEntriesLength() - 1)
					if nav.GetEntriesLength() < nav.NumVisibleLines {
						nav.SetStartLine(0)
					} else {
						nav.SetStartLine(nav.GetEntriesLength() - nav.NumVisibleLines)
					}
				}

			case pgDownKey:
				if nav.HasEntries() {
					terminal.ClearScreen()
					if nav.GetEntriesLength() < nav.GetCurrentIndex()+nav.NumVisibleLines {
						nav.SetCurrentIndex(nav.GetEntriesLength() - 1)
						nav.SetStartLine(nav.GetEntriesLength() - nav.NumVisibleLines)
					} else {
						nav.SetCurrentIndex(nav.GetCurrentIndex() + nav.NumVisibleLines)
						nav.SetStartLine(nav.GetCurrentIndex())
					}

					terminal.RenderOutput(&nav, &selected)
				}

			case pgUpKey:
				if nav.HasEntries() {
					if nav.GetCurrentIndex()-nav.NumVisibleLines < 0 {
						nav.SetCurrentIndex(0)
						nav.SetStartLine(0)
					} else {
						nav.SetCurrentIndex(nav.GetCurrentIndex() - nav.NumVisibleLines)
						nav.SetStartLine(nav.GetCurrentIndex())
					}

					terminal.RenderOutput(&nav, &selected)
				}

			case selectKey:
				if nav.HasEntries() {
					selected.DumpPrevious(nav.GetCurrentPath())
					selected.Select(nav.GetCurrentEntry())
					nav.IncrementIndex()
					terminal.RenderOutput(&nav, &selected)

				}

			case selectAllKey:
				if nav.HasEntries() {
					selected.DumpPrevious(nav.GetCurrentPath())
					selected.SelectAll(nav.GetEntries())
					terminal.RenderOutput(&nav, &selected)
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
						terminal.ClearScreen()
						task.Edit(nav.GetCurrentEntry().FullPath(), c.Cfg.DefaultEditor)
						terminal.RenderOutput(&nav, &selected)
					}
				}

			case renameKey:
				if nav.HasEntries() {
					task.WaitUserInput("rename", nav.GetCurrentEntry().Name, func(answ string) {
						if strings.Contains(answ, string(os.PathSeparator)) {
							utils.ShowErrAndContinue(fmt.Errorf("path separator can't be used inside name"))

							terminal.ResetFlushOutput(&nav, &selected)
							return
						}

						if answ == "" {
							utils.ShowErrAndContinue(fmt.Errorf("empty name?"))
							terminal.ResetFlushOutput(&nav, &selected)
							return
						}

						_, err := os.Stat(filepath.Join(nav.GetCurrentPath(), answ))
						if err == nil {

							utils.ShowErrAndContinue(fmt.Errorf("entry already exists"))
							terminal.ResetFlushOutput(&nav, &selected)
							return

						}

						err = os.Rename(nav.GetCurrentEntry().FullPath(), filepath.Join(nav.GetCurrentPath(), answ))

						if err != nil {
							panic(err)
						}

						nav.GetCurrentEntry().Name = answ
						terminal.ResetFlushOutput(&nav, &selected)
					})
				}

			case enterKey:
				if nav.HasEntries() {
					switch nav.GetCurrentEntry().IsDir {
					case false:
						task.ExecuteDefault(nav.GetCurrentEntry().FullPath(), utils.PrintDefaultErrorAndExit)
					default:
						enterSubfolder(&nav, &selected)
					}
				}
			case deleteKey, deleteKey2, deleteKey3:
				selected.DumpPrevious(nav.GetCurrentPath())
				if !selected.IsEmpty() {
					task.WaitUserInput("Delete selected entries", "y", func(answ string) {
						if answ == "y" {
							for key := range selected.GetAll() {
								err := os.RemoveAll(key.FullPath())
								if err != nil {
									if errors.Is(err, os.ErrPermission) {
										utils.ShowErrAndContinue(err)
										return
									} else {
										fmt.Println(err)
										os.Exit(1)
									}
								}
							}

							refresh()
							terminal.ResetFlushOutput(&nav, &selected)
						}
					})
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

				terminal.RenderOutput(&nav, &selected)

			case execKey:
				task.WaitUserInput("run command: ", "", func(input string) {
					terminal.ClearScreen()
					task.ExecCommand(input)
				})

				terminal.ResetFlushOutput(&nav, &selected)

			case copyKey, moveKey:

				if selected.IsEmpty() {
					continue
				}

				var ask string

				if event.Key == moveKey {
					ask = "Move"
				} else {
					ask = "Copy"
				}

				task.WaitUserInput(fmt.Sprintf("%s %s", ask, "selected into the current directory? :"), "y", func(s string) {
					var mv func(string, string, string)

					writeFile := func(srcPath, srcName, targetPath string) {
						newFilepath := filepath.Join(targetPath, srcName)

						srcFile, err := os.Open(filepath.Join(srcPath, srcName))
						if err != nil {
							panic(err)
						}

						defer srcFile.Close()

						targetFile, err := os.Create(newFilepath)
						if err != nil {
							panic(err)
						}
						defer targetFile.Close()

						_, err = io.Copy(targetFile, srcFile)
						if err != nil {
							panic(err)
						}

						info, err := os.Stat(filepath.Join(srcPath, srcName))
						if err != nil {
							fmt.Println(err)
							os.Exit(0)

						}

						err = os.Chmod(newFilepath, info.Mode())
						if err != nil {
							panic(err)
						}
					}

					mv = func(srcPath, srcName, targetDir string) {
						srcInfo, err := os.Stat(filepath.Join(srcPath, srcName))
						if err != nil {
							fmt.Println(err)
							os.Exit(0)
						}

						if srcInfo.IsDir() {
							proceed := true

							err := os.Mkdir(filepath.Join(targetDir, srcName), srcInfo.Mode())
							if err != nil {
								if errors.Is(err, os.ErrExist) {
									task.WaitUserInput(fmt.Sprintf("%s %s %s", "Folder ", srcName, " already exists, do you wish to merge them?"), "y", func(answ string) {
										// fmt.Println("Answer is:", answ)
										// time.Sleep(time.Millisecond * 500)

										if answ == "y" || answ == strings.ToLower("YES") {
											proceed = true
										} else {
											proceed = false
										}
									})
								}
							}

							if proceed {
								// fmt.Println("proceediung")
								// time.Sleep(time.Millisecond * 500)
								dc, err := os.ReadDir(filepath.Join(srcPath, srcName))
								if err != nil {
									fmt.Println(err)
									done <- true

								}

								for _, entry := range dc {
									mv(filepath.Join(srcPath, srcName), entry.Name(), filepath.Join(targetDir, srcName))
								}

							}

							// if key == moveKey {
							// 	err := os.RemoveAll(filepath.Join(srcPath, srcName))
							// 	if err != nil {
							// 		panic(err)
							// 	}
							// }

						} else {

							_, err := os.Stat(filepath.Join(targetDir, srcName))

							if err != nil && errors.Is(err, os.ErrNotExist) {
								// fmt.Println("File ", srcName, " not exisitng yet")
								// time.Sleep(time.Millisecond * 500)
								writeFile(srcPath, srcName, targetDir)
							} else {
								// fmt.Println("File ", srcName, " already exists")
								// time.Sleep(time.Millisecond * 500)
								task.WaitUserInput(fmt.Sprintf("%s %s %s", "File", srcName, " already exists, do you wish to overwrite it?"), "y", func(answ string) {
									if answ == "y" || answ == strings.ToLower("YES") {
										os.Remove(filepath.Join(targetDir, srcName))
										writeFile(srcPath, srcName, targetDir)
									}
								})
							}

						}
					}

					if s == "y" {
						for entry := range selected.GetAll() {

							if *entry.Path == nav.GetCurrentPath() {
								utils.ShowErrAndContinue(errors.New("copying / moving within same directory is not permitted"))
								return
							}

							if strings.HasPrefix(nav.GetCurrentPath(), entry.FullPath()) {
								utils.ShowErrAndContinue(errors.New("cannot move / copy a folder into itself"))
								return
							}

							mv(*entry.Path, entry.Name, nav.GetCurrentPath())

							refresh()

						}

						// TODO: after copying sort according to current sorting alghoritm
					}
				})

			case newFileKey:

				if task.CreateFsFile(nav.GetCurrentPath()) {
					//	sortMode = entry.SortByName(entries)

					refresh()
				}

			case newDirKey:

				if task.CreateFsDirectory(nav.GetCurrentPath()) {
					refresh()
				}

			case quitKey, quitKey2:
				terminal.ClearScreen()

				keyListener = false
				return

			}

		}

	}
}

func cleanup() {
	// Close keyboard, do any other necessary cleanup here
	terminal.ShowCursor()
	fmt.Println(terminal.ResetFmt)
	err := k.Close()
	if err != nil {
		fmt.Println("Error closing keyboard:", err)
	}
}
