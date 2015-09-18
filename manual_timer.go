package tme

import (
	"sync"
	"time"
)

// ManualTimer is a Timer implementation in which time must be manually advanced.
// Useful for writing deterministic tests against code that relies on real time.
type ManualTimer struct {
	sync.RWMutex
	ackFn      func()
	doneCh     chan Ack
	doneSig    chan struct{} // used to signal all waiters when doneTime is valid
	doneTime   time.Time
	stoppedSig chan struct{} // used to signal all waiters when stoppedTime is valid
	stopped    bool
}

// NewManualTimer creates a new ManualTimer
func NewManualTimer(ackFn func()) *ManualTimer {
	return &ManualTimer{
		ackFn:      ackFn,
		doneCh:     make(chan Ack),
		doneSig:    make(chan struct{}),
		doneTime:   time.Time{},
		stoppedSig: make(chan struct{}),
		stopped:    false,
	}
}

// MarkDone marks the timer done as if the timer expired. After this call, all channels
// returned by Done() in the past and future will receive and be closed. This func will
// block if no MarkDone calls have been made before and no other goroutines are listening on Chan()
//
// returns true if this is the first call to MarkDone, false otherwise
func (t *ManualTimer) MarkDone() bool {
	t.Lock()
	if t.doneTime.IsZero() {
		t.doneTime = time.Now()
		sig := t.doneSig
		t.Unlock()
		sig <- struct{}{}
		return true
	}
	t.Unlock()
	return false
}

// Stop is the interface implementation
func (t *ManualTimer) Stop() bool {
	t.Lock()
	defer t.Unlock()
	if t.stopped {
		return false
	}
	t.stopped = true
	close(t.stoppedSig)
	return true
}

// Done is the interface implementation.
func (t *ManualTimer) Done() <-chan Ack {
	t.RLock()
	defer t.RUnlock()
	go func() {
		select {
		case <-t.stoppedSig:
			return
		case <-t.doneSig:
			t.RLock()
			dt := t.doneTime
			t.RUnlock()
			t.doneCh <- Ack{Time: dt, Fn: t.ackFn}
		}
	}()
	return t.doneCh
}
