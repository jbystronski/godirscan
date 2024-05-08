package menu

import (
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
)

func QuitMenu() *MenuController {
	opts := []MenuOption{
		{
			Label: "Quit application",
			Event: event.QUIT_APP,
		},
		{
			Label: "Cancel",
			Event: event.Q,
		},
	}

	c := NewMenuController(opts, Dimensions{Width: 23, Height: 6})

	c.On(event.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:

			c.RunDefault(opt.Event)
		}
	})

	return c
}
