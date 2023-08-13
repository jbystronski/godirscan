package task

import "time"

var Ticker *time.Ticker

func StopTicker() {
	if Ticker != nil {
		Ticker.Stop()
	}
}

func StartTicker() {
	StopTicker()
	Ticker = time.NewTicker(time.Millisecond * 400)
}
