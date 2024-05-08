package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"sync"
	"time"

	"github.com/jbystronski/godirscan/pkg/app/data"
	"github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/lib/converter"
)

func (c *FsController) search() {
	var test func(fs.FileInfo) bool

	answ := c.getInput("Find by (1) pattern (2) size", "")

	if answ != "1" && answ != "2" {
		return
	}

	path := c.getInput("Path to search from", c.root)

	if path == "" {
		return
	}

	pattern := c.getInput("Search pattern", "")

	if answ == "1" {
		test = func(info fs.FileInfo) bool {
			return true
		}
	}

	if answ == "2" {

		answ := c.getInput("Select unit (0) bytes (1) kb (2) mb (3) gb", "")

		if answ != "0" && answ != "1" && answ != "2" && answ != "3" {
			return
		}

		unit, _ := strconv.Atoi(answ)

		unitName := converter.StorageUnits[unit]

		answ = c.getInput(fmt.Sprintf("%s %s", "Type min value in", unitName), "0")

		if answ == "" {
			return
		}

		min, err := strconv.ParseFloat(answ, 64)
		if err != nil {

			c.sendError(err)

			return
		}

		if min < 0 {
			min = 0
		}

		answ = c.getInput(fmt.Sprintf("%s %s", "Type max value ( no value means unlimited ) in ", unitName), "")

		if answ == "" {
			answ = "0"
		}

		max, err := strconv.ParseFloat(answ, 64)
		if err != nil {
			c.sendError(err)

			return
		}

		if max < 0 || min < 0 {
			c.sendError(errors.New("negative values are not allowed"))

			return
		}

		if max > 0 && max < min {

			c.sendError(fmt.Errorf("value: %v, can't be lower than min value: %v", max, min))

			return
		}

		minV, maxV := converter.ToBytes(unit, min, max)

		test = func(info fs.FileInfo) bool {
			if info.Size() >= minV {
				if maxV != 0 && info.Size() <= maxV || maxV == 0 {
					return true
				}
			}

			return false
		}

	}

	c.data.Reset()
	c.Navigator.Reset()

	_, tickerDone := c.ObserveTicker(c.ctx.Ctx.Done(), time.Millisecond*500, func() {
		c.render()
	})

	var selectWg sync.WaitGroup

	done := make(chan struct{}, 1)

	entryChan := make(chan struct {
		Path string
		Info fs.FileInfo
	})

	selectWg.Add(1)

	go func() {
		defer func() {
			selectWg.Done()
			tickerDone <- struct{}{}
		}()

		for {
			select {

			case <-done:

				return

			case v := <-entryChan:
				newEntry := &data.FsEntry{}
				newEntry.SetFsType(data.SearchResultDatatype)

				newEntry.SetName(v.Path)
				newEntry.SetSize(int(v.Info.Size()))
				c.data.Insert(newEntry)
			}
		}
	}()

	c.ctx.Observe(func() {
		err := global.Search(path, pattern, entryChan, test)
		if err != nil {
			c.sendError(err)
		}
	})

	done <- struct{}{}

	selectWg.Wait()
}
