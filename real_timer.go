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
	timer   *time.Timer
	stopped int32
}

func NewRealTimer(dur time.Duration) *RealTimer {
	timer := time.NewTimer(dur)
	return &RealTimer{timer: timer, stopped: 0}
}

func (r *RealTimer) Done() <-chan Ack {
	ch := make(chan Ack)
	go func() {
		for {
			if atomic.LoadInt32(&r.stopped) != 0 {
				return
			}
			recvT := <-r.timer.C
			ch <- Ack{Time: recvT, Fn: func() {}}
		}
	}()
	return ch
}

func (r *RealTimer) Stop() bool {
	if atomic.CompareAndSwapInt32(&r.stopped, 0, 1) {
		r.timer.Stop()
		return true
	}
	return false
}
