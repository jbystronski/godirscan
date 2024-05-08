package menu

import (
	"github.com/jbystronski/godirscan/pkg/app/config"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

func BookmarkGroupListMenu() *MenuController {
	opts := []MenuOption{
		{
			Label: "Main menu",

			Event: pubsub.M,
		},
		{
			Label: "Add bookmark group",

			Event: pubsub.BOOKMARK_ADD_GROUP,
		},
	}
	for groupName := range config.Running().Bookmarks {
		opts = append(opts, MenuOption{
			Label:       "Group",
			Description: groupName,

			Event: pubsub.BOOKMARK_GROUP,
		})
	}

	c := NewMenuController(opts, Dimensions{Height: 15, Width: termui.NewTerminal().Cols() / 3})

	c.On(pubsub.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {

		case true:
			if opt.Event == pubsub.BOOKMARK_GROUP {
				c.Publish("bookmark_group", pubsub.Message(opt.Description))
			}
			fallthrough
		default:
			c.RunDefault(opt.Event)

		}
	})

	return c
}
