package controller

import (
	"sort"
	"strconv"
)

func (c *Controller) sort() {
	answ := c.wrapInput("Sort by (1) Name (2) Size ASC (3) Size DESC (4) Type", "")

	if answ != "1" && answ != "2" && answ != "3" && answ != "4" {
		return
	}

	s, _ := strconv.ParseInt(answ, 0, 8)

	c.SetDefaultSort(s)

	switch s {
	case 1:
		c.DataAccessor.SortByName()

	case 2:
		sort.Sort(c.DataAccessor)

	case 3:
		sort.Sort(sort.Reverse(c.DataAccessor))

	case 4:
		c.DataAccessor.SortByType()

	default:
		c.DataAccessor.SortByName()

	}

	c.fullRender()
}

func (c *Controller) DefaultSort() {
	switch c.defaultSort {
	case 1:
		c.DataAccessor.SortByName()
	case 2:
		sort.Sort(c.DataAccessor)
	case 3:
		sort.Sort(sort.Reverse(c.DataAccessor))
	case 4:
		c.DataAccessor.SortByType()
	default:
		c.DataAccessor.SortByName()
	}
}

func (c *Controller) SetDefaultSort(s int64) {
	c.defaultSort = s
}
