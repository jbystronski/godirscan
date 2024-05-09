package menu

import (
	"github.com/jbystronski/godirscan/pkg/global/event"
)

func SortMenu() *MenuController {
	options := []MenuOption{
		{
			Label: "Main menu",

			Event: event.M,
		},
		{
			Label: "Sort by name",

			Event: event.SORT_NAME,
		},
		{
			Label: "Sort by type",

			Event: event.SORT_TYPE,
		},
		{
			Label: "Sort by size (ascending)",

			Event: event.SORT_SIZE_ASC,
		},
		{
			Label: "Sort by size (descending)",

			Event: event.SORT_SIZE_DESC,
		},
	}

	c := NewMenuController(options, Dimensions{Width: 40, Height: 7})

	c.On(event.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:
			c.RunDefault(opt.Event)
		}
	})

	return c
}
