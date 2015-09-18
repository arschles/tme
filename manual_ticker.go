package tme

import (
	"sync"
	"sync/atomic"
	"time"
)

// ManualTicker is a Ticker implementation that must be manually advanced.
// Useful for writing deterministic tests against code that relies on real time
type ManualTicker struct {
	sync.RWMutex
	ackFn   func()
	tickCh  chan Ack
	stopped int32
}

func NewManualTicker(ackFn func()) *ManualTicker {
	return &ManualTicker{
		ackFn:   ackFn,
		tickCh:  make(chan Ack),
		stopped: 0,
	}
}

func (t *ManualTicker) Chan() <-chan Ack {
	return t.tickCh
}

// Tick manually triggers a new tick. It's a synchronous tick, so the other side
// must receive on Chan() for this func to return
func (t *ManualTicker) Tick() {
	if atomic.LoadInt32(&t.stopped) == 1 {
		return
	}
	t.tickCh <- Ack{Time: time.Now(), Fn: t.ackFn}
}

func (t *ManualTicker) Stop() bool {
	return atomic.CompareAndSwapInt32(&t.stopped, 0, 1)
}
