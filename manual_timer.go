package tme

import (
	"sync/atomic"
	"time"
)

// ManualTimer is a Timer implementation in which time must be manually advanced.
// Useful for writing deterministic tests against code that relies on real time.
type ManualTimer struct {
	doneSig     chan struct{} // used to signal all waiters when doneTime is valid
	doneTime    int64
	stoppedSig  chan struct{} // used to signal all waiters when stoppedTime is valid
	stoppedTime int64
}

// NewManualTimer creates a new ManualTimer. only use this func to create ManualTimers
func NewManualTimer() *ManualTimer {
	return &ManualTimer{
		doneSig:     make(chan struct{}),
		doneTime:    0,
		stoppedSig:  make(chan struct{}),
		stoppedTime: 0,
	}
}

// MarkDone marks the timer done as if the timer expired. After this call, all channels
// returned by Done() in the past and future will receive and be closed
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

// Done is the interface implementation. All channels returned by Done prior to
// MarkDone being called will receive and close. All channels returned
func (t *ManualTimer) Done() <-chan Ack {
	ch := make(chan Ack)
	go func() {
		defer close(ch)
		<-t.doneSig
		doneTime := time.Unix(atomic.LoadInt64(&t.doneTime), 0)
		ch <- Ack{Time: doneTime, Fn: func() {}} //TODO: ack channel
	}()
	return ch
}
