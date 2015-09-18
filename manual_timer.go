package tme

import (
	"sync/atomic"
	"time"
)

// ManualTimer is a Timer implementation in which time must be manually advanced.
// Useful for writing deterministic tests against code that relies on real time.
type ManualTimer struct {
	ackFn       func()
	doneCh      chan Ack
	doneSig     chan struct{} // used to signal all waiters when doneTime is valid
	doneTime    int64
	stoppedSig  chan struct{} // used to signal all waiters when stoppedTime is valid
	stoppedTime int64
}

// NewManualTimer creates a new ManualTimer
func NewManualTimer(ackFn func()) *ManualTimer {
	return &ManualTimer{
		ackFn:       ackFn,
		doneCh:      make(chan Ack),
		doneSig:     make(chan struct{}),
		doneTime:    0,
		stoppedSig:  make(chan struct{}),
		stoppedTime: 0,
	}
}

// MarkDone marks the timer done as if the timer expired. After this call, all channels
// returned by Done() in the past and future will receive and be closed. This func will not
// block, and all calls but the first to MarkDone are noops
func (t *ManualTimer) MarkDone() {
	if atomic.CompareAndSwapInt64(&t.doneTime, 0, time.Now().Unix()) {
		close(t.doneSig)
	}
}

// Stop is the interface implementation
func (t *ManualTimer) Stop() bool {
	if atomic.CompareAndSwapInt64(&t.stoppedTime, 0, time.Now().Unix()) {
		close(t.stoppedSig)
		return true
	}
	return false
}

// Done is the interface implementation.
func (t *ManualTimer) Done() <-chan Ack {
	go func() {
		select {
		case <-t.stoppedSig:
			return
		case <-t.doneSig:
			doneTime := time.Unix(atomic.LoadInt64(&t.doneTime), 0)
			t.doneCh <- Ack{Time: doneTime, Fn: t.ackFn}
		}
	}()
	return t.doneCh
}
