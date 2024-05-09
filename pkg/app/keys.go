package app

import (
	"sync"

	"github.com/eiannone/keyboard"

	"github.com/jbystronski/godirscan/pkg/global/event"
	"github.com/jbystronski/pubsub"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

var (
	keyboardInit sync.Once
	n            *pubsub.Node
	k            *Keys
)

type Keys struct {
	*pubsub.Node
	keysMap map[keyboard.Key]pubsub.Event
	runeMap map[rune]pubsub.Event
}

func NewKeyboard() *pubsub.Node {
	keyboardInit.Do(func() {
		n = pubsub.NewNode(pubsub.GlobalBroker())

		k = &Keys{
			n,
			map[keyboard.Key]pubsub.Event{
				keyboard.KeyArrowDown:  event.ARROW_DOWN,
				keyboard.KeyArrowLeft:  event.ARROW_LEFT,
				keyboard.KeyArrowUp:    event.ARROW_UP,
				keyboard.KeyArrowRight: event.ARROW_RIGHT,
				keyboard.KeyEsc:        event.ESC,
				keyboard.KeyTab:        event.TAB,
				keyboard.KeyInsert:     event.INSERT,
				keyboard.KeyCtrlA:      event.CTRL_A,
				keyboard.KeyCtrlB:      event.CTRL_B,
				keyboard.KeyCtrlF:      event.CTRL_F,
				keyboard.KeyCtrlR:      event.CTRL_R,
				keyboard.KeyCtrlS:      event.CTRL_S,
				keyboard.KeyCtrlV:      event.CTRL_V,
				keyboard.KeyCtrlN:      event.CTRL_N,
				keyboard.KeyF7:         event.F7,
				keyboard.KeyF6:         event.F6,
				keyboard.KeyDelete:     event.DELETE,
				keyboard.KeyCtrlK:      event.CTRL_K,
				keyboard.KeyHome:       event.HOME,
				keyboard.KeyEnd:        event.END,
				keyboard.KeyPgdn:       event.PG_DOWN,
				keyboard.KeyPgup:       event.PG_UP,
				keyboard.KeyEnter:      event.ENTER,
			},
			map[rune]pubsub.Event{
				'h': event.H,
				'q': event.Q,
				'l': event.L,
				'i': event.I,
				'm': event.M,
				's': event.S,
				't': event.T,
				'd': event.D,
				'f': event.F,
				'e': event.E,
				'c': event.C,
				'g': event.G,
				'r': event.R,
			},
		}

		k.Init()
	})

	return n
}

func (keys *Keys) Init() {
	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		keys.Publish("err", pubsub.Message(err.Error()))
	}

	go func() {
		for {
			select {
			// case <-ctx.Done():
			// 	termui.ClearScreen()

			// 	return

			case ev := <-keysEvents:
				switch true {

				case termui.NewTerminal().IsCommandLineOpen():

					termui.NewTerminal().SendToCommandLine() <- struct {
						Key  keyboard.Key
						Char rune
					}{
						Key:  ev.Key,
						Char: ev.Rune,
					}

				case ev.Rune != 0:
					if str, ok := keys.runeMap[ev.Rune]; ok {
						keys.Passthrough(str, keys.Next())
					}

				case ev.Key != 0:
					if str, ok := keys.keysMap[ev.Key]; ok {
						keys.Passthrough(str, keys.Next())
					}
				}
			}
		}
	}()
}
