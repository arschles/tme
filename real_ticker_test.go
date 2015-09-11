package tme

import (
	"testing"
	"time"
)

func TestRealTickerSynchronousTicks(t *testing.T) {
	const dur = 1 * time.Millisecond
	const numChecks = 5
	ticker := NewRealTicker(dur)
	defer ticker.Stop()
	ch := ticker.Chan()
	for i := 0; i < numChecks; i++ {
		select {
		case <-ch:
		case <-time.After(dur * 2):
			t.Fatalf("ticker didn't tick within %s", dur*2)
		}
	}
}
