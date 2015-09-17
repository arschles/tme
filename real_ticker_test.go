package tme

import (
	"testing"
	"time"
)

func TestRealTickerSynchronousTicks(t *testing.T) {
	const dur = 1 * time.Millisecond
	const numChecks = 5
	ticker := NewRealTicker(dur, func() {})
	defer ticker.Stop()
	ch := ticker.Chan()
	for i := 0; i < numChecks; i++ {
		select {
		case <-ch:
		case <-time.After(dur * 2):
			t.Fatalf("ticker didn't tick within %s (thread %d)", dur*2, i)
		}
	}
}

func TestRealTickerStop(t *testing.T) {
	const dur = 10 * time.Millisecond
	ticker := NewRealTicker(dur, func() {})
	bch := ticker.Chan()
	select {
	case <-bch:
	case <-time.After(dur * 2):
		t.Errorf("ticker didn't tick within %s", dur*2)
	}
	if ticker.Stop() != true {
		t.Errorf("ticker was already marked stopped")
	}
	select {
	case <-bch:
		t.Errorf("ticker ticked after stop called")
	case <-time.After(dur * 2):
	}
	select {
	case <-ticker.Chan():
		t.Errorf("new channel from ticker ticked after stop called")
	case <-time.After(dur * 2):
	}
}
