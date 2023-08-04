package utils

import "time"

var ticker time.Ticker

func stopTicker(ticker *time.Ticker) {
	ticker.Stop()
}

func startTicker(ticker *time.Ticker) {
	stopTicker(ticker)
	*ticker = *time.NewTicker(time.Millisecond * 100)
}
