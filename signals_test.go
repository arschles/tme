package tme

import (
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	const dur = 20 * time.Millisecond
	chans := []chan struct{}{make(chan struct{}), make(chan struct{}), make(chan struct{})}
	s := newSignals()
	for _, ch := range chans {
		s.add(ch)
	}

	go func() {
		s.broadcast()
	}()

	for i, ch := range chans {
		select {
		case <-ch:
		case <-time.After(dur):
			t.Errorf("chan %d didn't receive within %s", i, dur)
		}
	}
}
