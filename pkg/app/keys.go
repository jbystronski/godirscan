package app

import (
	"sync"

	"github.com/eiannone/keyboard"

	e "github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/message"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

var (
	keyboardInit sync.Once
	n            *e.Node
	k            *Keys
)

type Keys struct {
	*e.Node
	keysMap map[keyboard.Key]e.Event
	runeMap map[rune]e.Event
}

func NewKeyboard() *e.Node {
	keyboardInit.Do(func() {
		n = e.NewNode()

		k = &Keys{
			n,
			map[keyboard.Key]e.Event{
				keyboard.KeyArrowDown:  e.ARROW_DOWN,
				keyboard.KeyArrowLeft:  e.ARROW_LEFT,
				keyboard.KeyArrowUp:    e.ARROW_UP,
				keyboard.KeyArrowRight: e.ARROW_RIGHT,
				keyboard.KeyEsc:        e.ESC,
				keyboard.KeyTab:        e.TAB,
				keyboard.KeyInsert:     e.INSERT,
				keyboard.KeyCtrlA:      e.CTRL_A,
				keyboard.KeyCtrlB:      e.CTRL_B,
				keyboard.KeyCtrlF:      e.CTRL_F,
				keyboard.KeyCtrlR:      e.CTRL_R,
				keyboard.KeyCtrlS:      e.CTRL_S,
				keyboard.KeyCtrlV:      e.CTRL_V,
				keyboard.KeyF7:         e.F7,
				keyboard.KeyF6:         e.F6,
				keyboard.KeyDelete:     e.DELETE,
				keyboard.KeyCtrlK:      e.CTRL_K,
				keyboard.KeyHome:       e.HOME,
				keyboard.KeyEnd:        e.END,
				keyboard.KeyPgdn:       e.PG_DOWN,
				keyboard.KeyPgup:       e.PG_UP,
				keyboard.KeyEnter:      e.ENTER,
			},
			map[rune]e.Event{
				'h': e.H,
				'q': e.Q,
				'l': e.L,
				'i': e.I,
				'm': e.M,
				's': e.S,
				't': e.T,
				'd': e.D,
				'f': e.F,
				'e': e.E,
				'c': e.C,
				'g': e.G,
			},
		}

		k.Init()
	})

	return n
}

func (keys *Keys) Init() {
	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		keys.Publish("err", message.Message(err.Error()))
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
						keys.Passthrough(str, keys.Next)
					}

				case ev.Key != 0:
					if str, ok := keys.keysMap[ev.Key]; ok {
						keys.Passthrough(str, keys.Next)
					}
				}
			}
		}
	}()
}
