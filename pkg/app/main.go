package app

import (
	"context"
	"sync"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/app/boxes"

	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/message"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
)

var (
	once sync.Once
	app  *event.Node
)

func New() *event.Node {
	defer func() {
		keyboard.Close()
		cls()
		showCursor()
	}()

	once.Do(func() {
		done := make(chan struct{}, 1)

		context.WithCancel(context.Background())

		keys := NewKeyboard()
		keys.Watch()
		resizer := NewResizeListener()

		resizer.Watch()

		start := NewStart()

		start.Watch()

		app := event.NewNode()

		keys.LinkTo(resizer).LinkTo(app).LinkTo(start)

		app.Subscribe("err", func(m message.Message) {
			errorScreen := boxes.NewError(string(m))
			errorScreen.Watch()

			app.Last().LinkTo(errorScreen)
			app.Passthrough(event.RENDER, app.Last())
		})

		app.OnGlobal(event.T, func() {
			config.Running().ChangeTheme(config.CurrentTheme)
		})

		app.OnGlobal(event.QUIT_APP, func() {
			app.UnlinkAll()
			cls()
			done <- struct{}{}
		})

		app.OnGlobal(event.RESIZE, func() {
			updateDimensions()

			if cols() >= MIN_WIDTH && rows() >= MIN_HEIGHT {
				if app.Next == boxes.NewResizeWarning() {
					app.LinkTo(start)
					app.Passthrough(event.RENDER, app.Next)

				}
			}

			if cols() < MIN_WIDTH || rows() < MIN_HEIGHT {
				cls()

				app.UnlinkAll()
				warn := boxes.NewResizeWarning()
				warn.Watch()
				app.LinkTo(warn)
				app.Passthrough(event.RENDER, app.Next)

			}
		})

		//	app.Passthrough(event.RENDER, app.Next)
		app.Watch()
		app.Passthrough(event.RESIZE, app)

		<-done
	})

	return app
}
