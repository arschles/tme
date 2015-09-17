package tme

import (
	"sync"
	"sync/atomic"
	"time"
)

// RealTicker is a Ticker implementation that uses real time
type RealTicker struct {
	sync.RWMutex
	ackFn    func()
	stopped  int32
	sigs     *signals
	tickTime time.Time
}

func NewRealTicker(dur time.Duration, ackFn func()) *RealTicker {
	t := &RealTicker{
		ackFn:    ackFn,
		stopped:  0,
		sigs:     newSignals(),
		tickTime: time.Time{},
	}
	go func() {
		ticker := time.NewTicker(dur)
		defer ticker.Stop()
		for {
			if atomic.LoadInt32(&t.stopped) == 1 {
				return
			}
			t.Lock()
			t.tickTime = <-ticker.C
			t.Unlock()
			t.sigs.broadcast()
		}
	}()
	return t
}

func (t *RealTicker) Chan() <-chan Ack {
	ch := make(chan Ack)
	sigCh := make(chan struct{})
	t.sigs.add(sigCh)
	go func() {
		for {
			if atomic.LoadInt32(&t.stopped) == 1 {
				return
			}
			<-sigCh
			t.RLock()
			tickTime := t.tickTime
			t.RUnlock()
			ch <- Ack{Time: tickTime, Fn: t.ackFn}
		}
	}()
	return ch
}

func (t *RealTicker) Stop() bool {
	return atomic.CompareAndSwapInt32(&t.stopped, 0, 1)
}
