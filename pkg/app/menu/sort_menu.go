package menu

import (
	e "github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
)

func SortMenu() *MenuController {
	options := []MenuOption{
		{
			Label: "Main menu",

			Event: e.M,
		},
		{
			Label: "Sort by name",

			Event: e.SORT_NAME,
		},
		{
			Label: "Sort by type",

			Event: e.SORT_TYPE,
		},
		{
			Label: "Sort by size (ascending)",

			Event: e.SORT_SIZE_ASC,
		},
		{
			Label: "Sort by size (descending)",

			Event: e.SORT_SIZE_DESC,
		},
	}

	c := NewMenuController(options, Dimensions{Width: 40, Height: 7})

	c.On(e.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:
			c.RunDefault(opt.Event)
		}
	})

	return c
}
