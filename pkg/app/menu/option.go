package menu

import "github.com/jbystronski/godirscan/pkg/lib/pubsub"

type MenuOption struct {
	Label, Description string

	Event pubsub.Event
}
