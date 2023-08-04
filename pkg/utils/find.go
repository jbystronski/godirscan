package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func promptFindByName(defaultPath string, entries []*entry.Entry, done chan<- struct{}) {
	var startPath, pattern string

	waitUserInput("Find (in path): ", defaultPath, func(s string) {
		startPath = s
		waitUserInput("Find (pattern): ", "", func(s string) {
			pattern = s

			startTicker(&ticker)

			entries = nil
			terminal.ClearScreen()

			go func() {
				findByName(startPath, pattern, &entries)
				done <- struct{}{}
			}()
		})
	})
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

func findBySize(allEntries *[]*entry.Entry, root, pattern string, min, max int64) {
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

							*allEntries = append(*allEntries, newEntry)

						}
					}
				}
			}
		}
	}

	find(filepath.Join(root), regexp.MustCompile(pattern), min, max)
}
