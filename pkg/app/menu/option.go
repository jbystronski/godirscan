package menu

import "github.com/jbystronski/pubsub"

type MenuOption struct {
	Label, Description string

	Event pubsub.Event
}
