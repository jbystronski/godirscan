package menu

import (
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
)

func QuitMenu() *MenuController {
	opts := []MenuOption{
		{
			Label: "Quit application",
			Event: pubsub.QUIT_APP,
		},
		{
			Label: "Cancel",
			Event: pubsub.Q,
		},
	}

	c := NewMenuController(opts, Dimensions{Width: 23, Height: 6})

	c.On(pubsub.ENTER, func() {
		opt := c.Options[c.Index()]

		switch true {
		default:

			c.RunDefault(opt.Event)
		}
	})

	return c
}
