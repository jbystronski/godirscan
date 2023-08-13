package listener

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
)

var resumeListener = make(chan struct{}, 1)

var stopListener = make(chan struct{}, 1)

var EventChan = make(chan keyboard.KeyEvent)

var exit = make(chan struct{}, 1)

var listenerOn bool

func Resume() {
	defer func() {
		fmt.Println("closing keyboard")
		err := keyboard.Close()
		if err != nil {
			panic(err)
		}
		// time.Sleep(time.Second * 1)
	}()

	fmt.Println("resuming listener")

	time.Sleep(time.Second * 1)

	// run := func() {
	// 	select {
	// 	case <-exit:
	// 		fmt.Println("exiting listener goroutine")
	// 		return
	// 	default:
	// 		keysEvents, err := keyboard.GetKeys(1)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		EventChan <- <-keysEvents

	// 	}
	// }

	for {
		select {

		// case <-resumeListener:

		// 	listenerOn = true

		// 	if !keyboard.IsStarted(time.Millisecond * 50) {
		// 		keyboard.Open()
		// 	}
		// 	go func() {
		// 		run()
		// 	}()

		case <-stopListener:

			fmt.Println("stopping listener")

			time.Sleep(time.Second * 1)

			return
		default:
			keysEvents, err := keyboard.GetKeys(1)
			if err != nil {
				panic(err)
			}
			EventChan <- <-keysEvents

		}
	}
}

func Stop() {
	stopListener <- struct{}{}
}

// func Resume() {
// 	resumeListener <- struct{}{}
// }
