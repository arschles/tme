package tme

import (
	"testing"
	"time"
)

func TestRealTickerSynchronousTicks(t *testing.T) {
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
	case <-time.After(dur):
	}
	select {
	case <-ticker.Chan():
		t.Errorf("new channel from ticker ticked after stop called")
	case <-time.After(dur):
	}
}

func TestRealTickerAck(t *testing.T) {
	ackCh := make(chan struct{})
	ticker := NewRealTicker(dur, func() { ackCh <- struct{}{} })
	defer ticker.Stop()
	select {
	case ack := <-ticker.Chan():
		go func() { ack.Fn() }()
	case <-time.After(dur * 2):
		t.Errorf("ticker didn't tick within %s", dur*2)
	}

	select {
	case <-ackCh:
	case <-time.After(dur):
		t.Errorf("ack didn't happen within %s", dur)
	}
}

func TestRealTickerMultipleRecv(t *testing.T) {
	const n = 10
	ticker := NewRealTicker(dur, func() {})
	defer ticker.Stop()
	numRecvs := numTickerRecv(ticker, 10)
	if numRecvs != 1 {
		t.Errorf("%d ticker receives", numRecvs)
	}
}
