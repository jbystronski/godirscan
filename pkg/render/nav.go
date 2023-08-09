package nav

// package render

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"math"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"time"

// 	k "github.com/eiannone/keyboard"
// 	"github.com/jbystronski/godirscan/pkg/cache"
// 	"github.com/jbystronski/godirscan/pkg/config"
// 	"github.com/jbystronski/godirscan/pkg/converter"
// 	"github.com/jbystronski/godirscan/pkg/entry"
// 	"github.com/jbystronski/godirscan/pkg/render"
// 	"github.com/jbystronski/godirscan/pkg/task"
// 	"github.com/jbystronski/godirscan/pkg/terminal"
// )

// func GetParentPath(path string) string {
// 	parentDir, _ := filepath.Split(path)

// 	if parentDir == string(os.PathSeparator) {
// 		return parentDir
// 	}

// 	return strings.TrimSuffix(parentDir, string(os.PathSeparator))
// }

// func pathExists(dir, file string) bool {
// 	_, err := os.Stat(filepath.Join(dir, file))

// 	if err != nil && errors.Is(err, os.ErrNotExist) {
// 		return false
// 	}

// 	return true
// }

// func RunMainLoop() {
// 	go func() {
// 		render.PrintBanner()
// 		k.Open()
// 		currentPath, entries = task.ScanInputDirectory(config.Cfg.DefaultRootDirectory)
// 		k.Close()
// 		rootPath = currentPath
// 		entry.SetSort(&sortMode, entries)

// 		render.ResetFlushOutput()

// 		go func() {
// 			task.ScanDirectorySize(&entries, &dirSize)

// 			sizeCalculationDone <- struct{}{}
// 		}()

// 		for {
// 			task.StartTicker()
// 			keysEvents, err := k.GetKeys(1)
// 			if err != nil {
// 				panic(err)
// 			}

// 			select {
// 			case <-sizeCalculationDone:

// 				task.StopTicker()

// 				cache.Store(currentPath, dirSize, entries)

// 				render.RenderOutput()

// 			case <-task.Ticker.C:
// 				render.RenderOutput()

// 			case <-searchDone:

// 				if len(entries) == 0 {
// 					fmt.Println("no entries found")
// 				}

// 			case event := <-keysEvents:

// 				switch event.Key {

// 				case menuKey:

// 					render.ClearScreen()

// 					render.PrintHelp()

// 					render.ResetFlushOutput()

// 				case k.KeyCtrlC:
// 					os.Exit(0)

// 				case scanKey:

// 					newRootDir, newEntries := task.ScanInputDirectory(config.Cfg.DefaultRootDirectory)

// 					if newRootDir != "" {

// 						selected = make(map[*entry.Entry]struct{})
// 						cache.Clear()
// 						currentPath = newRootDir

// 						entries = newEntries

// 						render.ResetFlushOutput()
// 						task.StartTicker()

// 						go func() {
// 							task.ScanDirectorySize(&entries, &dirSize)

// 							sizeCalculationDone <- struct{}{}
// 						}()

// 					}

// 				case downKey:
// 					if moveDown() {
// 						renderOutput()
// 					}

// 				case upKey:
// 					if moveUp() {
// 						renderOutput()
// 					}

// 				case rightKey:

// 					enterSubfolder()

// 				case leftKey:

// 					if currentPath != rootPath {

// 						path := GetParentPath(currentPath)

// 						if cachedEntries, ok := cache.Get(path); ok {

// 							entries = cachedEntries.Entries
// 							dirSize = cachedEntries.Size
// 							currentPath = path

// 							startLine = 0
// 							endLine = 0
// 							currentIndex = backTrace[len(backTrace)-1]
// 							sortMode = 0
// 							entry.SetSort(&sortMode, entries)
// 							flushOutput()

// 						} else {

// 							task.StartTicker()
// 							currentPath, entries = task.ScanDirectory(path)

// 							startLine = 0
// 							endLine = 0
// 							currentIndex = backTrace[len(backTrace)-1]
// 							sortMode = 0
// 							entry.SetSort(&sortMode, entries)
// 							flushOutput()

// 							go func() {
// 								task.ScanDirectorySize(&entries, &dirSize)
// 								sizeCalculationDone <- struct{}{}
// 							}()

// 						}

// 					}

// 				case findKey:

// 					task.PromptFindByName(currentPath, entries, searchDone)

// 					// var startPath, pattern string

// 					// waitUserInput("Find (in path): ", currentPath, func(s string) {
// 					// 	startPath = s
// 					// 	waitUserInput("Find (pattern): ", "", func(s string) {
// 					// 		pattern = s

// 					// 		startTicker(&ticker)

// 					// 		entries = nil
// 					// 		terminal.ClearScreen()

// 					// 		go func() {
// 					// 			findByName(startPath, pattern, &entries)
// 					// 			searchDone <- struct{}{}
// 					// 		}()
// 					// 	})
// 					// })

// 				case findSizeKey:

// 					var unitAsString, pathName, pattern string
// 					var unitAsInt int
// 					var min, max float64 = 0, math.MaxFloat64

// 					waitUserInput("Find by size, unit: ( 0=bytes 1=kb 2=mb 3=gb ) ", "2", func(s string) {
// 						number, _ := strconv.Atoi(s)

// 						if number < 0 || number > len(converter.StorageUnits)-1 {
// 							showErrAndContinue(errors.New("invalid unix index"))
// 							return
// 						}

// 						unitAsString = converter.StorageUnits[number]
// 						unitAsInt = number

// 						waitUserInput(fmt.Sprintf("%s %s", "Type min value in", unitAsString), "0", func(s string) {
// 							num, err := strconv.ParseFloat(s, 64)
// 							if err != nil {
// 								showErrAndContinue(err)
// 								return
// 							}

// 							if num < 0 {
// 								num = 0
// 							}
// 							min = num

// 							waitUserInput(fmt.Sprintf("%s %s", "Type max value ( 0 or no value means unlimited ) in", unitAsString), "", func(s string) {
// 								if s != "" || s != "0" {
// 									n, err := strconv.ParseFloat(s, 64)
// 									if err != nil {
// 										showErrAndContinue(err)
// 										return
// 									}
// 									max = n

// 								}

// 								if max < min {
// 									showErrAndContinue(fmt.Errorf("max value: %v, can't be lower than min value: %v", max, min))
// 									return
// 								}

// 								waitUserInput("Root directory to search from: ", "", func(s string) {
// 									_, err := os.Stat(s)
// 									if err != nil {
// 										showErrAndContinue(err)
// 										return
// 									}

// 									pathName = s

// 									waitUserInput("Pattern to match: ", "", func(s string) {
// 										pattern = s

// 										minV, maxV := converter.ToBytes(converter.StorageUnits[unitAsInt], min, max)

// 										startTicker(&ticker)

// 										entries = nil
// 										terminal.ClearScreen()

// 										go func() {
// 											findBySize(&entries, pathName, pattern, minV, maxV)
// 											searchDone <- struct{}{}
// 										}()
// 									})
// 								})
// 							})
// 						})
// 					})

// 				case sortKey:
// 					if len(entries) > 0 {
// 						entry.SetSort(&sortMode, entries)
// 						renderOutput()
// 					}

// 				case homeKey:
// 					if len(entries) > 0 {
// 						currentIndex = 0
// 						startLine = 0
// 						renderOutput()
// 					}

// 				case endKey:
// 					if len(entries) > 0 {
// 						currentIndex = len(entries) - 1
// 						if len(entries) < numVisibleLines {
// 							startLine = 0
// 						} else {
// 							startLine = len(entries) - numVisibleLines
// 						}

// 						renderOutput()
// 					}

// 				case pgDownKey:

// 					if len(entries) > 0 {
// 						terminal.ClearScreen()

// 						if len(entries) < currentIndex+numVisibleLines {
// 							currentIndex = len(entries) - 1
// 							startLine = len(entries) - numVisibleLines

// 						} else {
// 							currentIndex = currentIndex + numVisibleLines
// 							startLine = currentIndex
// 						}

// 						renderOutput()
// 					}

// 				case pgUpKey:

// 					if len(entries) > 0 {
// 						if currentIndex-numVisibleLines < 0 {

// 							currentIndex = 0
// 							startLine = 0

// 						} else {
// 							currentIndex = currentIndex - numVisibleLines
// 							startLine = currentIndex
// 						}

// 						renderOutput()

// 					}

// 				case selectKey, selectKey2:

// 					if len(entries) > 0 {
// 						clearSelected()

// 						selectEntry(entries[currentIndex])

// 						incrementCursor()
// 						renderOutput()
// 					}

// 				case selectAllKey:

// 					if len(entries) > 0 {
// 						clearSelected()

// 						selectAll()
// 						renderOutput()
// 					}

// 				case viewKey:

// 					if len(entries) > 0 {
// 						en := entries[currentIndex]

// 						if !en.IsDir {

// 							clearScreen()
// 							task.Peek(filepath.Join(currentPath, en.Name))
// 						}
// 					}

// 				case editKey:

// 					if len(entries) > 0 {

// 						en := entries[currentIndex]

// 						if !en.IsDir {
// 							terminal.ClearScreen()
// 							task.Edit(filepath.Join(currentPath, en.Name), Cfg.DefaultEditor, printDefaultErrorAndExit)

// 							renderOutput()

// 						}
// 					}

// 				case renameKey:

// 					if len(entries) > 0 {

// 						en := entries[currentIndex]

// 						waitUserInput("rename ", en.Name, func(answ string) {
// 							answ = strings.TrimSpace(answ)
// 							if strings.ContainsAny(answ, string(os.PathSeparator)) {
// 								showErrAndContinue(fmt.Errorf("path separator can't be used inside name"))
// 								return
// 							}

// 							if answ == "" {
// 								showErrAndContinue(fmt.Errorf("empty name?"))
// 								return
// 							}

// 							_, err := os.Stat(filepath.Join(currentPath, answ))
// 							if err == nil {

// 								showErrAndContinue(fmt.Errorf("entry already exists"))
// 								return

// 							}
// 							err = os.Rename(en.FullPath(), filepath.Join(currentPath, answ))
// 							if err != nil {
// 								panic(err)
// 							}

// 							en.Name = answ
// 							resetFlushOutput()
// 						})

// 					}

// 				case enterKey:

// 					if len(entries) > 0 {
// 						en := entries[currentIndex]
// 						switch en.IsDir {
// 						case false:

// 							executeDefault(filepath.Join(currentPath, en.Name), showErrAndContinue)

// 						default:

// 							enterSubfolder()

// 							// TODO: implement filesystemScan on path and change to scanMode
// 						}
// 					}

// 				case execKey:

// 					waitUserInput("run command: ", "", func(s string) {
// 						terminal.ClearScreen()
// 						execCommand(s, printDefaultErrorAndExit)
// 					})

// 					resetFlushOutput()

// 				case deleteKey, deleteKey2, deleteKey3:
// 					clearSelected()
// 					if len(selected) > 0 {
// 						waitUserInput("Delete selected entries", "y", func(answ string) {
// 							if answ == "y" {
// 								for key := range selected {
// 									err := os.RemoveAll(key.FullPath())
// 									if err != nil {
// 										if errors.Is(err, os.ErrPermission) {
// 											showErrAndContinue(err)
// 										} else {
// 											fmt.Println(err)
// 											os.Exit(1)
// 										}
// 									}
// 								}

// 								refresh()
// 								resetFlushOutput()
// 							}
// 						})
// 					}

// 				case nextThemeKey:
// 					{
// 						num := Cfg.CurrentSchema
// 						if num < uint(len(Cfg.ColorSchemas)-1) {
// 							num++
// 						} else {
// 							num = 0
// 						}

// 						switchTheme(num)
// 						renderOutput()
// 					}

// 				case newFileKey:

// 					if mode == scanMode && createFsFile(currentPath, &entries, showErrAndContinue) {
// 						sortMode = entry.SortByName(entries)

// 						renderOutput()
// 					}

// 				case newDirKey:

// 					if mode == scanMode && createFsDirectory(currentPath, &entries, showErrAndContinue) {
// 						sortMode = entry.SortByName(entries)
// 						renderOutput()
// 					}

// 				case copyKey, moveKey:

// 					if len(selected) == 0 {
// 						continue
// 					}

// 					var ask string

// 					if event.Key == moveKey {
// 						ask = "Move"
// 					} else {
// 						ask = "Copy"
// 					}

// 					waitUserInput(fmt.Sprintf("%s %s", ask, "selected into the current directory? :"), "y", func(s string) {
// 						var mv func(string, string, string)

// 						writeFile := func(srcPath, srcName, targetPath string) {
// 							newFilepath := filepath.Join(targetPath, srcName)

// 							srcFile, err := os.Open(filepath.Join(srcPath, srcName))
// 							if err != nil {
// 								panic(err)
// 							}

// 							defer srcFile.Close()

// 							targetFile, err := os.Create(newFilepath)
// 							if err != nil {
// 								panic(err)
// 							}
// 							defer targetFile.Close()

// 							_, err = io.Copy(targetFile, srcFile)
// 							if err != nil {
// 								panic(err)
// 							}

// 							info, err := os.Stat(filepath.Join(srcPath, srcName))
// 							if err != nil {
// 								fmt.Println(err)
// 								os.Exit(0)

// 							}

// 							err = os.Chmod(newFilepath, info.Mode())
// 							if err != nil {
// 								panic(err)
// 							}
// 						}

// 						mv = func(srcPath, srcName, targetDir string) {
// 							srcInfo, err := os.Stat(filepath.Join(srcPath, srcName))
// 							if err != nil {
// 								fmt.Println(err)
// 								os.Exit(0)
// 							}

// 							if srcInfo.IsDir() {
// 								proceed := true

// 								err := os.Mkdir(filepath.Join(targetDir, srcName), srcInfo.Mode())
// 								if err != nil {
// 									if errors.Is(err, os.ErrExist) {
// 										waitUserInput(fmt.Sprintf("%s %s %s", "Folder ", srcName, " already exists, do you wish to merge them?"), "y", func(answ string) {
// 											// fmt.Println("Answer is:", answ)
// 											// time.Sleep(time.Millisecond * 500)

// 											if answ == "y" || answ == strings.ToLower("YES") {
// 												proceed = true
// 											} else {
// 												proceed = false
// 											}
// 										})
// 									}
// 								}

// 								if proceed {
// 									// fmt.Println("proceediung")
// 									// time.Sleep(time.Millisecond * 500)
// 									dc, err := os.ReadDir(filepath.Join(srcPath, srcName))
// 									if err != nil {
// 										printDefaultErrorAndExit(err)
// 									}

// 									for _, entry := range dc {
// 										mv(filepath.Join(srcPath, srcName), entry.Name(), filepath.Join(targetDir, srcName))
// 									}

// 								}

// 								// if key == moveKey {
// 								// 	err := os.RemoveAll(filepath.Join(srcPath, srcName))
// 								// 	if err != nil {
// 								// 		panic(err)
// 								// 	}
// 								// }

// 							} else {

// 								_, err := os.Stat(filepath.Join(targetDir, srcName))

// 								if err != nil && errors.Is(err, os.ErrNotExist) {
// 									// fmt.Println("File ", srcName, " not exisitng yet")
// 									// time.Sleep(time.Millisecond * 500)
// 									writeFile(srcPath, srcName, targetDir)
// 								} else {
// 									// fmt.Println("File ", srcName, " already exists")
// 									// time.Sleep(time.Millisecond * 500)
// 									waitUserInput(fmt.Sprintf("%s %s %s", "File", srcName, " already exists, do you wish to overwrite it?"), "y", func(answ string) {
// 										if answ == "y" || answ == strings.ToLower("YES") {
// 											os.Remove(filepath.Join(targetDir, srcName))
// 											writeFile(srcPath, srcName, targetDir)
// 										}
// 									})
// 								}

// 							}
// 						}

// 						if s == "y" {
// 							for entry := range selected {

// 								if *entry.Path == currentPath {
// 									showErrAndContinue(errors.New("copying / moving within same directory is not permitted"))
// 									return
// 								}

// 								if strings.HasPrefix(currentPath, filepath.Join(*entry.Path, entry.Name)) {
// 									showErrAndContinue(errors.New("cannot move / copy a folder into itself"))
// 									return
// 								}

// 								mv(*entry.Path, entry.Name, currentPath)

// 								refresh()

// 							}

// 							// TODO: after copying sort according to current sorting alghoritm
// 						}
// 					})

// 				case quitKey, quitKey2:

// 					terminal.ClearScreen()

// 					// keyboardOpen = false
// 					fmt.Print("\033[?25h")
// 					done <- true

// 					time.Sleep(100 * time.Millisecond)
// 					os.Exit(0)
// 				}

// 			}
// 		}
// 	}()

// 	<-done
// }
