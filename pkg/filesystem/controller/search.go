package controller

import (
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"sync"
	"time"

	"github.com/jbystronski/godirscan/pkg/converter"
	"github.com/jbystronski/godirscan/pkg/filesystem"
	"github.com/jbystronski/godirscan/pkg/filesystem/utils"
)

func (c *Controller) search() {
	var test func(fs.FileInfo) bool

	answ := c.wrapInput("Find by (1) pattern (2) size", "")

	if answ != "1" && answ != "2" {
		return
	}
	path := c.wrapInput("Path to search from", c.path)

	if path == "" {
		return
	}

	pattern := c.wrapInput("Search pattern", "")

	if answ == "1" {
		test = func(info fs.FileInfo) bool {
			return true
		}
	}

	if answ == "2" {
		answ := c.wrapInput("Select unit (0) bytes (1) kb (2) mb (3) gb", "")

		if answ != "0" && answ != "1" && answ != "2" && answ != "3" {
			return
		}

		unit, _ := strconv.Atoi(answ)

		unitName := converter.StorageUnits[unit]

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

	c.DataAccessor.Reset()
	c.ChunkNavigator.Reset()

	_, tickerDone := c.ObserveTicker(c.ctx.Ctx.Done(), time.Millisecond*500, func() {
		c.SetChunk(c.DataAccessor.Len())
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
				newEntry := &filesystem.FsEntry{}
				newEntry.SetFsType(filesystem.SearchResult)

				newEntry.SetName(v.Path)
				newEntry.SetSize(int(v.Info.Size()))
				c.Insert(newEntry)
			}
		}
	}()

	c.ctx.Observe(func() {
		utils.Search(path, pattern, c.ErrorChan, entryChan, test)
	})

	done <- struct{}{}

	selectWg.Wait()
}
