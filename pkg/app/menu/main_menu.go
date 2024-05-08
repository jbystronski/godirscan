package menu

import (
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

func MainMenu() *MenuController {
	options := []MenuOption{
		{
			Label:       "Switch panel",
			Description: "(TAB)",
			Event:       event.TAB,
		},

		{
			Label:       "Select",
			Description: "(INSERT)",
			Event:       event.INSERT,
		},

		{
			Label:       "Enter / execute file",
			Description: "(ENTER)",
			Event:       event.ENTER,
		},

		{
			Label:       "Select / deselect all",
			Description: "(CTRL A)",
			Event:       event.CTRL_A,
		},
		{
			Label:       "Edit file",
			Description: "(E)",
			Event:       event.E,
		},
		{
			Label:       "Display file info",
			Description: "(I)",
			Event:       event.I,
		},
		{
			Label:       "Scan directory",
			Description: "(S)",
			Event:       event.S,
		},
		{
			Label:       "Rename file / directory",
			Description: "(CTRL R)",
			Event:       event.CTRL_R,
		},
		{
			Label:       "Create new directory",
			Description: "(D)",
			Event:       event.D,
		},
		{
			Label:       "Create new file",
			Description: "(F)",
			Event:       event.F,
		},
		{
			Label:       "Search",
			Description: "(CTRL F)",
			Event:       event.CTRL_F,
		},
		{
			Label:       "Copy selected",
			Description: "(CTRL V)",
			Event:       event.CTRL_V,
		},
		{
			Label:       "Move selected",
			Description: "(F6)",
			Event:       event.F6,
		},
		{
			Label:       "Go to entry",
			Description: "(G)",
			Event:       event.G,
		},
		{
			Label:       "Delete selected entries",
			Description: "(DELETE)",
			Event:       event.DELETE,
		},

		{
			Label:       "Change theme",
			Description: "(T)",
			Event:       event.T,
		},

		{
			Label:       "Close menu",
			Description: "(Q)",
			Event:       event.Q,
		},

		{
			Label: "Sort",
			Event: event.SORT_MENU,
		},
		{
			Label: "Edit settings (default location)",
			Event: event.SETTINGS,
		},
		{
			Label: "Bookmarks",
			Event: event.BOOKMARK_GROUP_LIST,
		},
		{
			Label: "Quit",
			Event: event.ESC,
		},
	}

	c := NewMenuController(options, Dimensions{Height: 15, Width: termui.NewTerminal().Cols() / 3})

	c.On(event.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:

			c.RunDefault(opt.Event)
		}
	})

	return c
}
