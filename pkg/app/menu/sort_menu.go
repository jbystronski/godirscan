package menu

import "github.com/jbystronski/godirscan/pkg/lib/pubsub"

func SortMenu() *MenuController {
	options := []MenuOption{
		{
			Label: "Main menu",

			Event: pubsub.M,
		},
		{
			Label: "Sort by name",

			Event: pubsub.SORT_NAME,
		},
		{
			Label: "Sort by type",

			Event: pubsub.SORT_TYPE,
		},
		{
			Label: "Sort by size (ascending)",

			Event: pubsub.SORT_SIZE_ASC,
		},
		{
			Label: "Sort by size (descending)",

			Event: pubsub.SORT_SIZE_DESC,
		},
	}

	c := NewMenuController(options, Dimensions{Width: 40, Height: 7})

	c.On(pubsub.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:
			c.RunDefault(opt.Event)
		}
	})

	return c
}
