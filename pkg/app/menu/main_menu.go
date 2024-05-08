package menu

import (
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

func MainMenu() *MenuController {
	options := []MenuOption{
		{
			Label:       "Switch panel",
			Description: "(TAB)",
			Event:       pubsub.TAB,
		},

		{
			Label:       "Select",
			Description: "(INSERT)",
			Event:       pubsub.INSERT,
		},

		{
			Label:       "Enter / execute file",
			Description: "(ENTER)",
			Event:       pubsub.ENTER,
		},

		{
			Label:       "Select / deselect all",
			Description: "(CTRL A)",
			Event:       pubsub.CTRL_A,
		},
		{
			Label:       "Edit file",
			Description: "(E)",
			Event:       pubsub.E,
		},
		{
			Label:       "Display file info",
			Description: "(I)",
			Event:       pubsub.I,
		},
		{
			Label:       "Scan directory",
			Description: "(S)",
			Event:       pubsub.S,
		},
		{
			Label:       "Rename file / directory",
			Description: "(CTRL R)",
			Event:       pubsub.CTRL_R,
		},
		{
			Label:       "Create new directory",
			Description: "(D)",
			Event:       pubsub.D,
		},
		{
			Label:       "Create new file",
			Description: "(F)",
			Event:       pubsub.F,
		},
		{
			Label:       "Search",
			Description: "(CTRL F)",
			Event:       pubsub.CTRL_F,
		},
		{
			Label:       "Copy selected",
			Description: "(CTRL V)",
			Event:       pubsub.CTRL_V,
		},
		{
			Label:       "Move selected",
			Description: "(F6)",
			Event:       pubsub.F6,
		},
		{
			Label:       "Go to entry",
			Description: "(G)",
			Event:       pubsub.G,
		},
		{
			Label:       "Delete selected entries",
			Description: "(DELETE)",
			Event:       pubsub.DELETE,
		},

		{
			Label:       "Change theme",
			Description: "(T)",
			Event:       pubsub.T,
		},

		{
			Label:       "Close menu",
			Description: "(Q)",
			Event:       pubsub.Q,
		},

		{
			Label: "Sort",
			Event: pubsub.SORT_MENU,
		},
		{
			Label: "Edit settings (default location)",
			Event: pubsub.SETTINGS,
		},
		{
			Label: "Bookmarks",
			Event: pubsub.BOOKMARK_GROUP_LIST,
		},
		{
			Label: "Quit",
			Event: pubsub.ESC,
		},
	}

	c := NewMenuController(options, Dimensions{Height: 15, Width: termui.NewTerminal().Cols() / 3})

	c.On(pubsub.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:

			c.RunDefault(opt.Event)
		}
	})

	return c
}
