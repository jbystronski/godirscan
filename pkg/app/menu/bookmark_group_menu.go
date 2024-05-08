package menu

import (
	"github.com/jbystronski/godirscan/pkg/app/config"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

func BookmarkGroupMenu(groupName string) *MenuController {
	opts := []MenuOption{
		{
			Label: "Main menu",

			Event: pubsub.M,
		},
		{
			Label: "All groups",

			Event: pubsub.BOOKMARK_GROUP_LIST,
		},
		{
			Label:       "Delete group",
			Description: groupName,
			Event:       pubsub.BOOKMARK_REMOVE_GROUP,
		},
		{
			Label:       "Add bookmark to",
			Description: groupName,
			Event:       pubsub.BOOKMARK,
		},
	}
	for _, bookmark := range config.Running().Bookmarks[groupName] {
		opts = append(opts, MenuOption{
			Label:       "Open",
			Description: bookmark,
			Event:       pubsub.BOOKMARK_OPEN,
		})

		opts = append(opts, MenuOption{
			Label:       "Delete",
			Description: bookmark,
			Event:       pubsub.BOOKMARK_REMOVE,
		})

	}

	c := NewMenuController(opts, Dimensions{Height: 15, Width: termui.NewTerminal().Cols() / 3})

	c.On(pubsub.ARROW_LEFT, func() {
		c.Unlink()
		c.Passthrough(pubsub.RENDER, c.Prev)
		c.Passthrough(pubsub.BOOKMARK_GROUP_LIST, c.Prev)
	})

	c.On(pubsub.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {

		case true:
			if opt.Event == pubsub.BOOKMARK {
				c.Publish("bookmark_group", pubsub.Message(groupName))
			}
			fallthrough
		case true:

			if opt.Event == pubsub.BOOKMARK_OPEN {
				c.Publish("bookmark", pubsub.Message(opt.Description))
			}
			fallthrough
		case true:
			if opt.Event == pubsub.BOOKMARK_REMOVE {
				c.Publish("bookmark_group", pubsub.Message(groupName))
				c.Publish("bookmark", pubsub.Message(opt.Description))
			}
			fallthrough

		case true:
			if opt.Event == pubsub.BOOKMARK_REMOVE_GROUP {
				c.Publish("bookmark_group", pubsub.Message(groupName))
			}
			fallthrough

		default:
			c.RunDefault(opt.Event)

		}
	})

	return c
}
