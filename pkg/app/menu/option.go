package menu

import "github.com/jbystronski/godirscan/pkg/lib/pubsub/event"

type MenuOption struct {
	Label, Description string

	Event event.Event
}
