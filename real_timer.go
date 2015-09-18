package tme

import (
	"sync/atomic"
	"time"
)

// RealTimer is a Timer implementation that uses real time. Example usage:
//
//  t := NewRealTimer(1 * time.Second)
//  go func() {
//    <-t.Done()
//    log.Printf("timer done")
//  }()
type RealTimer struct {
	ackFn   func()
	timer   *time.Timer
	stopped int32
	doneCh  chan Ack
}

func NewRealTimer(dur time.Duration, ackFn func()) *RealTimer {
	timer := time.NewTimer(dur)
	return &RealTimer{
		ackFn:   ackFn,
		timer:   timer,
		stopped: 0,
		doneCh:  make(chan Ack),
	}
}

func (r *RealTimer) Done() <-chan Ack {
	go func() {
		for {
			if atomic.LoadInt32(&r.stopped) != 0 {
				return
			}
			recvT := <-r.timer.C
			r.doneCh <- Ack{Time: recvT, Fn: r.ackFn}
		}
	}()
	return r.doneCh
}

func (r *RealTimer) Stop() bool {
	if atomic.CompareAndSwapInt32(&r.stopped, 0, 1) {
		r.timer.Stop()
		return true
	}
	return false
}
