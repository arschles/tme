package tme

import (
	"sync/atomic"
	"time"
)

// RealTicker is a Ticker implementation that uses real time
type RealTicker struct {
	ticker  *time.Ticker
	stopped int32
}

func NewRealTicker(dur time.Duration) *RealTicker {
	return &RealTicker{ticker: time.NewTicker(dur), stopped: 0}
}

func (t *RealTicker) Chan() <-chan Ack {
	ch := make(chan Ack)
	if atomic.LoadInt32(&t.stopped) == 1 {
		close(ch)
		return ch
	}
	go func() {
		defer close(ch)
		for atomic.LoadInt32(&t.stopped) == 0 {
			tickTime := <-t.ticker.C
			ch <- Ack{Time: tickTime, Fn: func() {}}
		}
	}()
	return ch
}

func (t *RealTicker) Stop() bool {
	return atomic.CompareAndSwapInt32(&t.stopped, 0, 1)
}
