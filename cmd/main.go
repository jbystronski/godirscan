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
	"time"

	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/cache"
	"github.com/jbystronski/godirscan/pkg/converter"
	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/terminal"
	"github.com/jbystronski/godirscan/pkg/ui"
)

var ticker time.Ticker

func startTicker(ticker *time.Ticker) {
	ticker.Stop()
	*ticker = *time.NewTicker(100 * time.Millisecond)
}

func stopTicker(ticker *time.Ticker) {
	ticker.Stop()
}

var spaceKey = k.KeySpace
var menuKey = k.KeyF2
var scanKey = k.KeyF9
var downKey = k.KeyArrowDown
var upKey = k.KeyArrowUp
var rightKey = k.KeyArrowRight
var leftKey = k.KeyArrowLeft
var findKey = k.KeyCtrlF
var findSizeKey = k.KeyCtrlL
var homeKey = k.KeyHome
var endKey = k.KeyEnd
var pgDownKey = k.KeyPgdn
var pgUpKey = k.KeyPgup
var sortKey = k.KeyCtrlS
var selectKey = k.KeyInsert
var selectAllKey = k.KeyCtrlA
var viewKey = k.KeyF3
var editKey = k.KeyCtrlU
var renameKey = k.KeyCtrlR
var enterKey = k.KeyEnter
var execKey = k.KeyCtrlE
var deleteKey = k.KeyF8
var deleteKey2 = k.KeyDelete
var nextThemeKey = k.KeyCtrlSlash
var newFileKey = k.KeyCtrlW
var newDirKey = k.KeyF7
var copyKey = k.KeyF5
var moveKey = k.KeyF6
var quitKey = k.KeyEsc

var quitKey2 = k.KeyF10
var quitChar = 'q'

func main() {
	const (
		scanMode = iota
		filterMode
		reservedLines = 2
	)

	var (
		selected = make(map[*entry.Entry]struct{})

		backTrace           = []int{}
		sizeCalculationDone = make(chan struct{})
		searchDone          = make(chan struct{})

		done = make(chan bool, 1)

		sortMode uint8 = 0
		entries  []*entry.Entry

		currentPath string
		rootPath    string

		selectedEntriesPath string

		mode            = scanMode
		numVisibleLines = terminal.GetNumVisibleLines() - reservedLines
		// numVisibleCols                   = terminal.GetNumVisibleCols()

		// done                                      = make(chan struct{})
		dirSize, currentIndex, startLine, endLine int
	)

	renderOutput := func() {
		terminal.MoveCursorTop()

		endLine = startLine + numVisibleLines

		if endLine > len(entries) {
			endLine = len(entries)
		}

		var sep string
		terminal.ClearLine()
		printHeader(currentPath + terminal.Space + entry.PrintSizeAsString(dirSize))

		if len(entries) == 0 {
			printEmpty()
		}
		for i := startLine; i < endLine; i++ {

			if i == len(entries)-1 {
				sep = terminal.CornerLine + terminal.Hseparator
			} else {
				sep = terminal.TeeLine + terminal.Hseparator
			}
			terminal.ClearLine()
			if currentIndex == i {
				highlightRow(sep, *entries[i])
			} else if _, ok := selected[entries[i]]; ok {
				markRow(sep, *entries[i])
			} else {
				printRow(sep, *entries[i])
			}

		}
		fmt.Print("\033[?25h") // Show cursor
	}

	flushOutput := func() {
		terminal.ClearScreen()
		renderOutput()
	}

	resetFlushOutput := func() {
		terminal.ClearScreen()
		currentIndex = 0
		startLine = 0
		endLine = 0
		renderOutput()
	}

	showErrAndContinue := func(err error) {
		if err != nil {
			// blockNavigation = true
			terminal.ClearScreen()
			fmt.Printf("%s%s\n\n%s%s", terminal.White, err, terminal.ResetFmt, "press space to continue")

			for {

				char, key, err := k.GetKey()
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}

				if key == spaceKey || char == quitChar {

					// blockNavigation = false
					resetFlushOutput()
					break

				}

			}

		}
	}

	incrementCursor := func() {
		if currentIndex < len(entries)-1 {
			currentIndex++
		}
	}

	clearSelected := func() {
		if currentPath != selectedEntriesPath {
			selected = make(map[*entry.Entry]struct{})
			selectedEntriesPath = currentPath

		}
	}

	refresh := func() {
		selected = make(map[*entry.Entry]struct{})

		cache.Clear()
		startTicker(&ticker)
		currentPath, entries = utils.ScanDirectory(currentPath)

		go func() {
			scanDirectorySize(&entries, &dirSize)
			sizeCalculationDone <- struct{}{}
		}()
	}

	enterSubfolder := func() {
		if mode == scanMode && len(entries) > 0 {

			en := entries[currentIndex]

			if en.IsDir && pathExists(currentPath, en.Name) {
				backTrace = append(backTrace, currentIndex)

				p := filepath.Join(currentPath, en.Name)

				if cachedEntries, ok := cache.Get(p); ok {
					entries = cachedEntries.Entries
					dirSize = cachedEntries.Size
					currentPath = p
					resetFlushOutput()
				} else {
					startTicker(&ticker)
					currentPath, entries = utils.ScanDirectory(p)
					sortMode = 0
					entry.SetSort(&sortMode, entries)
					resetFlushOutput()

					go func() {
						scanDirectorySize(&entries, &dirSize)
						sizeCalculationDone <- struct{}{}
					}()

				}

			}

		}
	}

	go func() {
		ui.PrintBanner()
		k.Open()
		currentPath, entries = scanInputDirectory(Cfg.DefaultRootDirectory, showErrAndContinue)
		k.Close()
		rootPath = currentPath
		entry.SetSort(&sortMode, entries)
		resetFlushOutput()

		go func() {
			scanDirectorySize(&entries, &dirSize)

			sizeCalculationDone <- struct{}{}
		}()

		for {
			startTicker(&ticker)
			keysEvents, err := k.GetKeys(1)
			if err != nil {
				panic(err)
			}

			select {
			case <-sizeCalculationDone:

				stopTicker(&ticker)

				cache.Store(currentPath, dirSize, entries)

				renderOutput()

			case <-ticker.C:
				renderOutput()

			case <-searchDone:

				if len(entries) == 0 {
					fmt.Println("no entries found")
				}

			case event := <-keysEvents:

				switch event.Key {

				case menuKey:

					terminal.ClearScreen()

					printHelp()

					resetFlushOutput()

				case scanKey:

					newRootDir, newEntries := scanInputDirectory(Cfg.DefaultRootDirectory, showErrAndContinue)

					if newRootDir != "" {

						mode = scanMode
						selected = make(map[*entry.Entry]struct{})
						cache.Clear()
						currentPath = newRootDir

						entries = newEntries

						resetFlushOutput()
						startTicker(&ticker)

						go func() {
							scanDirectorySize(&entries, &dirSize)

							sizeCalculationDone <- struct{}{}
						}()

					}

				case downKey:

					if len(entries) > 0 && currentIndex < len(entries)-1 {
						currentIndex++
						if currentIndex >= startLine+numVisibleLines {
							startLine++
						}
						renderOutput()
					}

				case upKey:

					if len(entries) > 0 && currentIndex > 0 {
						currentIndex--
						if currentIndex < startLine {
							startLine--
						}
						renderOutput()
					}

				case rightKey:

					enterSubfolder()

				case leftKey:

					if mode == scanMode {
						if currentPath != rootPath {

							path := GetParentPath(currentPath)

							if cachedEntries, ok := cache.Get(path); ok {

								entries = cachedEntries.Entries
								dirSize = cachedEntries.Size
								currentPath = path

								startLine = 0
								endLine = 0
								currentIndex = backTrace[len(backTrace)-1]
								sortMode = 0
								entry.SetSort(&sortMode, entries)
								flushOutput()

							} else {

								startTicker(&ticker)
								currentPath, entries = utils.ScanDirectory(path)

								startLine = 0
								endLine = 0
								currentIndex = backTrace[len(backTrace)-1]
								sortMode = 0
								entry.SetSort(&sortMode, entries)
								flushOutput()

								go func() {
									scanDirectorySize(&entries, &dirSize)
									sizeCalculationDone <- struct{}{}
								}()

							}

						}
					}

				case findKey:

					promptFindByName(currentPath, entries, searchDone)

					// var startPath, pattern string

					// waitUserInput("Find (in path): ", currentPath, func(s string) {
					// 	startPath = s
					// 	waitUserInput("Find (pattern): ", "", func(s string) {
					// 		pattern = s

					// 		startTicker(&ticker)

					// 		entries = nil
					// 		terminal.ClearScreen()

					// 		go func() {
					// 			findByName(startPath, pattern, &entries)
					// 			searchDone <- struct{}{}
					// 		}()
					// 	})
					// })

				case findSizeKey:

					var unitAsString, pathName, pattern string
					var unitAsInt int
					var min, max float64 = 0, math.MaxFloat64

					waitUserInput("Find by size, unit: ( 0=bytes 1=kb 2=mb 3=gb ) ", "2", func(s string) {
						number, _ := strconv.Atoi(s)

						if number < 0 || number > len(converter.StorageUnits)-1 {
							showErrAndContinue(errors.New("invalid unix index"))
							return
						}

						unitAsString = converter.StorageUnits[number]
						unitAsInt = number

						waitUserInput(fmt.Sprintf("%s %s", "Type min value in", unitAsString), "0", func(s string) {
							num, err := strconv.ParseFloat(s, 64)
							if err != nil {
								showErrAndContinue(err)
								return
							}

							if num < 0 {
								num = 0
							}
							min = num

							waitUserInput(fmt.Sprintf("%s %s", "Type max value ( 0 or no value means unlimited ) in", unitAsString), "", func(s string) {
								if s != "" || s != "0" {
									n, err := strconv.ParseFloat(s, 64)
									if err != nil {
										showErrAndContinue(err)
										return
									}
									max = n

								}

								if max < min {
									showErrAndContinue(fmt.Errorf("max value: %v, can't be lower than min value: %v", max, min))
									return
								}

								waitUserInput("Root directory to search from: ", "", func(s string) {
									_, err := os.Stat(s)
									if err != nil {
										showErrAndContinue(err)
										return
									}

									pathName = s

									waitUserInput("Pattern to match: ", "", func(s string) {
										pattern = s

										minV, maxV := converter.ToBytes(converter.StorageUnits[unitAsInt], min, max)

										startTicker(&ticker)

										entries = nil
										terminal.ClearScreen()

										go func() {
											findBySize(&entries, pathName, pattern, minV, maxV)
											searchDone <- struct{}{}
										}()
									})
								})
							})
						})
					})

				case sortKey:
					if len(entries) > 0 {
						entry.SetSort(&sortMode, entries)
						renderOutput()
					}

				case homeKey:
					if len(entries) > 0 {
						currentIndex = 0
						startLine = 0
						renderOutput()
					}

				case endKey:
					if len(entries) > 0 {
						currentIndex = len(entries) - 1
						if len(entries) < numVisibleLines {
							startLine = 0
						} else {
							startLine = len(entries) - numVisibleLines
						}

						renderOutput()
					}

				case pgDownKey:

					if len(entries) > 0 {
						terminal.ClearScreen()

						if len(entries) < currentIndex+numVisibleLines {
							currentIndex = len(entries) - 1
							startLine = len(entries) - numVisibleLines

						} else {
							currentIndex = currentIndex + numVisibleLines
							startLine = currentIndex
						}

						renderOutput()
					}

				case pgUpKey:

					if len(entries) > 0 {
						if currentIndex-numVisibleLines < 0 {

							currentIndex = 0
							startLine = 0

						} else {
							currentIndex = currentIndex - numVisibleLines
							startLine = currentIndex
						}

						renderOutput()

					}

				case selectKey:

					if len(entries) > 0 {
						clearSelected()

						selectEntry(entries[currentIndex], &selected)

						incrementCursor()
						renderOutput()
					}

				case selectAllKey:

					if len(entries) > 0 {
						clearSelected()

						selectAllEntries(entries, &selected)
						renderOutput()
					}

				case viewKey:

					if len(entries) > 0 {
						en := entries[currentIndex]

						if !en.IsDir {

							terminal.ClearScreen()
							peek(filepath.Join(currentPath, en.Name))
						}
					}

				case editKey:

					if len(entries) > 0 {

						en := entries[currentIndex]

						if !en.IsDir {
							terminal.ClearScreen()
							edit(filepath.Join(currentPath, en.Name), Cfg.DefaultEditor, printDefaultErrorAndExit)

							renderOutput()

						}
					}

				case renameKey:

					if len(entries) > 0 {

						en := entries[currentIndex]

						waitUserInput("rename ", en.Name, func(answ string) {
							answ = strings.TrimSpace(answ)
							if strings.ContainsAny(answ, string(os.PathSeparator)) {
								showErrAndContinue(fmt.Errorf("path separator can't be used inside name"))
								return
							}

							if answ == "" {
								showErrAndContinue(fmt.Errorf("empty name?"))
								return
							}

							_, err := os.Stat(filepath.Join(currentPath, answ))
							if err == nil {

								showErrAndContinue(fmt.Errorf("entry already exists"))
								return

							}
							err = os.Rename(en.FullPath(), filepath.Join(currentPath, answ))
							if err != nil {
								panic(err)
							}

							en.Name = answ
							resetFlushOutput()
						})

					}

				case enterKey:

					if len(entries) > 0 {
						en := entries[currentIndex]
						switch en.IsDir {
						case false:

							executeDefault(filepath.Join(currentPath, en.Name), showErrAndContinue)

						default:

							enterSubfolder()

							// TODO: implement filesystemScan on path and change to scanMode
						}
					}

				case execKey:

					waitUserInput("run command: ", "", func(s string) {
						terminal.ClearScreen()
						execCommand(s, printDefaultErrorAndExit)
					})

					resetFlushOutput()

				case deleteKey, deleteKey2:
					clearSelected()
					if len(selected) > 0 {
						waitUserInput("Delete selected entries", "y", func(answ string) {
							if answ == "y" {
								for key := range selected {
									err := os.RemoveAll(key.FullPath())
									if err != nil {
										if errors.Is(err, os.ErrPermission) {
											showErrAndContinue(err)
										} else {
											fmt.Println(err)
											os.Exit(1)
										}
									}
								}

								refresh()
								resetFlushOutput()
							}
						})
					}

				case nextThemeKey:
					{
						num := Cfg.CurrentSchema
						if num < uint(len(Cfg.ColorSchemas)-1) {
							num++
						} else {
							num = 0
						}

						switchTheme(num)
						renderOutput()
					}

				case newFileKey:

					if mode == scanMode && createFsFile(currentPath, &entries, showErrAndContinue) {
						sortMode = entry.SortByName(entries)

						renderOutput()
					}

				case newDirKey:

					if mode == scanMode && createFsDirectory(currentPath, &entries, showErrAndContinue) {
						sortMode = entry.SortByName(entries)
						renderOutput()
					}

				case copyKey, moveKey:

					if len(selected) == 0 {
						continue
					}

					var ask string

					if event.Key == moveKey {
						ask = "Move"
					} else {
						ask = "Copy"
					}

					waitUserInput(fmt.Sprintf("%s %s", ask, "selected into the current directory? :"), "y", func(s string) {
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
										waitUserInput(fmt.Sprintf("%s %s %s", "Folder ", srcName, " already exists, do you wish to merge them?"), "y", func(answ string) {
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
										printDefaultErrorAndExit(err)
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
									waitUserInput(fmt.Sprintf("%s %s %s", "File", srcName, " already exists, do you wish to overwrite it?"), "y", func(answ string) {
										if answ == "y" || answ == strings.ToLower("YES") {
											os.Remove(filepath.Join(targetDir, srcName))
											writeFile(srcPath, srcName, targetDir)
										}
									})
								}

							}
						}

						if s == "y" {
							for entry := range selected {

								if *entry.Path == currentPath {
									showErrAndContinue(errors.New("copying / moving within same directory is not permitted"))
									return
								}

								if strings.HasPrefix(currentPath, filepath.Join(*entry.Path, entry.Name)) {
									showErrAndContinue(errors.New("cannot move / copy a folder into itself"))
									return
								}

								mv(*entry.Path, entry.Name, currentPath)

								refresh()

							}

							// TODO: after copying sort according to current sorting alghoritm
						}
					})

				case quitKey, quitKey2:

					terminal.ClearScreen()

					// keyboardOpen = false
					fmt.Print("\033[?25h")
					done <- true

					time.Sleep(100 * time.Millisecond)
					os.Exit(0)
				}

			}
		}
	}()

	<-done
}
