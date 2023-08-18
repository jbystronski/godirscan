package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func PromptFindByName(defaultPath string, entries []*entry.Entry, done chan<- struct{}) {
	startPath, err := WaitInput("Find (in path): ", defaultPath)
	if err != nil {
		panic(err)
	}

	pattern, err := WaitInput("Find (pattern): ", "")
	if err != nil {
		panic(err)
	}

	StartTicker()

	entries = nil
	terminal.ClearScreen()

	go func() {
		findByName(startPath, pattern, &entries)
		done <- struct{}{}
	}()
}

func findByName(root, pattern string, allEntries *[]*entry.Entry) {
	var find func(*regexp.Regexp, string)

	find = func(reg *regexp.Regexp, path string) {
		dc, err := os.ReadDir(path)
		if err != nil {
			fmt.Println(err)
			// showErrAndContinue(err)
			return
		}

		for _, dirEntry := range dc {

			f, _ := dirEntry.Info()

			if f.IsDir() {
				find(reg, filepath.Join(path, f.Name()))
			} else {
				if reg.Match([]byte(f.Name())) {
					*allEntries = append(*allEntries, &entry.Entry{
						Name:  filepath.Join(path, f.Name()),
						Size:  int(f.Size()),
						IsDir: f.IsDir(),
						Path:  &path,
					})
				}
			}
		}
	}

	find(regexp.MustCompile(pattern), filepath.Join(root))
}

// func PromptFindBySize() {
// 	var pathName string
// 	var unitAsInt int
// 	var max float64 = 0, math.MaxFloat64

// 	unitAsString := WaitInput("Find by size, unit: ( 0=bytes 1=kb 2=mb 3=gb ) ", "2")

// 	unit, _ := strconv.Atoi(unitAsString)

// 	if number < 0 || number > len(converter.StorageUnits)-1 {
// 		utils.ShowErrAndContinue(errors.New("invalid unit index"))
// 		return
// 	}

// 	unitAsString = converter.StorageUnits[number]
// 	unitAsInt = number

// 	min := task.WaitInput(fmt.Sprintf("%s %s", "Type min value in", unitAsString), "0")

// 	num, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		utils.ShowErrAndContinue(err)
// 		return
// 	}

// 	if num < 0 {
// 		num = 0
// 	}
// 	min = num

// 	max = task.WaitInput(fmt.Sprintf("%s %s", "Type max value ( 0 or no value means unlimited ) in", unitAsString), "")
// 	if max != "" || max != "0" {
// 		n, err := strconv.ParseFloat(max, 64)
// 		if err != nil {
// 			utils.ShowErrAndContinue(err)
// 			return
// 		}
// 		max = n

// 	}

// 	if max < min {
// 		utils.ShowErrAndContinue(fmt.Errorf("max value: %v, can't be lower than min value: %v", max, min))
// 		return
// 	}

// 	dir := task.WaitInput("Root directory to search from: ", "")
// 	_, err = os.Stat(dir)
// 	if err != nil {
// 		utils.ShowErrAndContinue(err)
// 		return
// 	}

// 	pattern := task.WaitInput("Pattern to match: ", "")

// 	minV, maxV := converter.ToBytes(converter.StorageUnits[unitAsInt], min, max)

// 	task.StartTicker()

// 	nav.ClearEntries()
// 	terminal.ClearScreen()
// }

func findBySize(allEntries []*entry.Entry, root, pattern string, min, max int64) {
	var find func(string, *regexp.Regexp, int64, int64)

	find = func(path string, reg *regexp.Regexp, min, max int64) {
		dc, err := os.ReadDir(path)
		if err != nil {
			fmt.Println(err)
			// showErrAndContinue(err)
			return
		}

		for _, dirEntry := range dc {

			f, _ := dirEntry.Info()

			if f.IsDir() {
				find(filepath.Join(path, f.Name()), reg, min, max)
			} else {
				if f.Size() >= min {
					if max != 0 && f.Size() <= max || max == 0 {
						if reg.Match([]byte(f.Name())) {

							newEntry := &entry.Entry{
								Name:  filepath.Join(path, f.Name()),
								Size:  int(f.Size()),
								IsDir: dirEntry.IsDir(),
								Path:  &path,
							}

							allEntries = append(allEntries, newEntry)

						}
					}
				}
			}
		}
	}

	find(filepath.Join(root), regexp.MustCompile(pattern), min, max)
}
