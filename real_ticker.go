package tme

import (
	"sync/atomic"
	"time"
)

// RealTicker is a Ticker implementation that uses real time
type RealTicker struct {
	ackFn   func()
	tickCh  chan Ack
	stopped int32
	stopCh  chan struct{}
}

func NewRealTicker(dur time.Duration, ackFn func()) *RealTicker {
	t := &RealTicker{
		ackFn:   ackFn,
		tickCh:  make(chan Ack),
		stopped: 0,
		stopCh:  make(chan struct{}),
	}
	go func() {
		ticker := time.NewTicker(dur)
		defer ticker.Stop()
		for {
			if atomic.LoadInt32(&t.stopped) == 1 {
				return
			}
			select {
			case tickTime := <-ticker.C:
				t.tickCh <- Ack{Time: tickTime, Fn: t.ackFn}
			case <-t.stopCh:
				return
			}
		}
	}()
	return t
}

func (t *RealTicker) Chan() <-chan Ack {
	return t.tickCh
}

func (t *RealTicker) Stop() bool {
	if atomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		close(t.stopCh)
		return true
	}
	return false
}
