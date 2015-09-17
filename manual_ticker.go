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
	ackFn    func()
	sigs     *signals
	tickTime time.Time
	stopped  int32
}

func NewManualTicker(ackFn func()) *ManualTicker {
	return &ManualTicker{
		ackFn:    ackFn,
		sigs:     newSignals(),
		tickTime: time.Time{},
		stopped:  0,
	}
}

func (t *ManualTicker) Chan() <-chan Ack {
	ch := make(chan Ack)
	if atomic.LoadInt32(&t.stopped) == 1 {
		return ch
	}

	sigCh := make(chan struct{})
	t.sigs.add(sigCh)
	go func() {
		<-sigCh
		t.RLock()
		ch <- Ack{Time: t.tickTime, Fn: t.ackFn}
		t.RUnlock()
	}()
	return ch
}

// Tick manually triggers a new tick
func (t *ManualTicker) Tick() {
	if atomic.LoadInt32(&t.stopped) == 1 {
		return
	}
	t.sigs.broadcast()
}

func (t *ManualTicker) Stop() bool {
	return atomic.CompareAndSwapInt32(&t.stopped, 0, 1)
}
