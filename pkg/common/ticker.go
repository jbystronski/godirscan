package common

import (
	"time"
)

type Ticker struct {
	ticker   time.Ticker
	interval time.Duration
	init     bool
}

func (t *Ticker) Start(duration time.Duration) {
	t.Stop()
	t.ticker = *time.NewTicker(duration)
	t.init = true
}

func (t *Ticker) Stop() {
	if t.init {

		t.ticker.Stop()
		t.init = false

	}
}

func (t *Ticker) Interval() time.Duration {
	return t.interval
}

func (t *Ticker) SetInterval(i time.Duration) {
	t.interval = i
}

func (t *Ticker) Tick() <-chan time.Time {
	return t.ticker.C
}

func (t *Ticker) IsInitialized() bool {
	return t.init
}
