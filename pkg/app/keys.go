package app

import (
	"sync"

	"github.com/eiannone/keyboard"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"

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
		n = pubsub.NewNode()

		k = &Keys{
			n,
			map[keyboard.Key]pubsub.Event{
				keyboard.KeyArrowDown:  pubsub.ARROW_DOWN,
				keyboard.KeyArrowLeft:  pubsub.ARROW_LEFT,
				keyboard.KeyArrowUp:    pubsub.ARROW_UP,
				keyboard.KeyArrowRight: pubsub.ARROW_RIGHT,
				keyboard.KeyEsc:        pubsub.ESC,
				keyboard.KeyTab:        pubsub.TAB,
				keyboard.KeyInsert:     pubsub.INSERT,
				keyboard.KeyCtrlA:      pubsub.CTRL_A,
				keyboard.KeyCtrlB:      pubsub.CTRL_B,
				keyboard.KeyCtrlF:      pubsub.CTRL_F,
				keyboard.KeyCtrlR:      pubsub.CTRL_R,
				keyboard.KeyCtrlS:      pubsub.CTRL_S,
				keyboard.KeyCtrlV:      pubsub.CTRL_V,
				keyboard.KeyF7:         pubsub.F7,
				keyboard.KeyF6:         pubsub.F6,
				keyboard.KeyDelete:     pubsub.DELETE,
				keyboard.KeyCtrlK:      pubsub.CTRL_K,
				keyboard.KeyHome:       pubsub.HOME,
				keyboard.KeyEnd:        pubsub.END,
				keyboard.KeyPgdn:       pubsub.PG_DOWN,
				keyboard.KeyPgup:       pubsub.PG_UP,
				keyboard.KeyEnter:      pubsub.ENTER,
			},
			map[rune]pubsub.Event{
				'h': pubsub.H,
				'q': pubsub.Q,
				'l': pubsub.L,
				'i': pubsub.I,
				'm': pubsub.M,
				's': pubsub.S,
				't': pubsub.T,
				'd': pubsub.D,
				'f': pubsub.F,
				'e': pubsub.E,
				'c': pubsub.C,
				'g': pubsub.G,
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
