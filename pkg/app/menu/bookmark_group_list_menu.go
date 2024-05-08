package menu

import (
	"github.com/jbystronski/godirscan/pkg/app/config"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/message"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

func BookmarkGroupListMenu() *MenuController {
	opts := []MenuOption{
		{
			Label: "Main menu",

			Event: event.M,
		},
		{
			Label: "Add bookmark group",

			Event: event.BOOKMARK_ADD_GROUP,
		},
	}
	for groupName := range config.Running().Bookmarks {
		opts = append(opts, MenuOption{
			Label:       "Group",
			Description: groupName,

			Event: event.BOOKMARK_GROUP,
		})
	}

	c := NewMenuController(opts, Dimensions{Height: 15, Width: termui.NewTerminal().Cols() / 3})

	c.On(event.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {

		case true:
			if opt.Event == event.BOOKMARK_GROUP {
				c.Publish("bookmark_group", message.Message(opt.Description))
			}
			fallthrough
		default:
			c.RunDefault(opt.Event)

		}
	})

	return c
}
