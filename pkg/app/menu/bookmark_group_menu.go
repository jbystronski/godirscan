package menu

import (
	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/global/event"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
	"github.com/jbystronski/pubsub"
)

func BookmarkGroupMenu(groupName string) *MenuController {
	opts := []MenuOption{
		{
			Label: "Main menu",

			Event: event.M,
		},
		{
			Label: "All groups",

			Event: event.BOOKMARK_GROUP_LIST,
		},
		{
			Label:       "Delete group",
			Description: groupName,
			Event:       event.BOOKMARK_REMOVE_GROUP,
		},
		{
			Label:       "Add bookmark to",
			Description: groupName,
			Event:       event.BOOKMARK,
		},
	}
	for _, bookmark := range config.Running().Bookmarks[groupName] {
		opts = append(opts, MenuOption{
			Label:       "Open",
			Description: bookmark,
			Event:       event.BOOKMARK_OPEN,
		})

		opts = append(opts, MenuOption{
			Label:       "Delete",
			Description: bookmark,
			Event:       event.BOOKMARK_REMOVE,
		})

	}

	c := NewMenuController(opts, Dimensions{Height: 15, Width: termui.NewTerminal().Cols() / 3})

	c.On(event.ARROW_LEFT, func() {
		c.Unlink()
		c.Passthrough(event.RENDER, c.Prev())
		c.Passthrough(event.BOOKMARK_GROUP_LIST, c.Prev())
	})

	c.On(event.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {

		case true:
			if opt.Event == event.BOOKMARK {
				c.Publish("bookmark_group", pubsub.Message(groupName))
			}
			fallthrough
		case true:

			if opt.Event == event.BOOKMARK_OPEN {
				c.Publish("bookmark", pubsub.Message(opt.Description))
			}
			fallthrough
		case true:
			if opt.Event == event.BOOKMARK_REMOVE {
				c.Publish("bookmark_group", pubsub.Message(groupName))
				c.Publish("bookmark", pubsub.Message(opt.Description))
			}
			fallthrough

		case true:
			if opt.Event == event.BOOKMARK_REMOVE_GROUP {
				c.Publish("bookmark_group", pubsub.Message(groupName))
			}
			fallthrough

		default:
			c.RunDefault(opt.Event)

		}
	})

	return c
}
