package task

import "time"

var Ticker time.Ticker

func StopTicker() {
	Ticker.Stop()
}

func StartTicker() {
	Ticker.Stop()
	Ticker = *time.NewTicker(time.Millisecond * 100)
}
